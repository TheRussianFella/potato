package slave

import (
	"encoding/json"
	"net"
	"time"
)

//////////
// Communication with a client
//////////

// StartServing begins an infinite loop for serving connections.
func (s *PotatoSlave) StartServing() {

	listener, err := net.Listen("tcp4", ":"+s.port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for i := s.numToServ; i != 0; i-- {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// Check if there are workers available
		select {
		case <-s.availableWorkers:

			name, _ := s.authConnection(c)
			go s.handleConnection(c, name)

		case <-time.After(time.Second):

			json.NewEncoder(c).Encode(ResponseMessage{
				Code:          _NW,
				StatusMessage: statusMessages[_NW],
				Value:         "",
			})
		}
	}
}

// CommandMessage is a structure that describes command messages sent by a client
// to a slave node
type CommandMessage struct {
	Name      string
	Arguments []string
	TTL       time.Duration
}

// ResponseMessage is a message sent back to user
type ResponseMessage struct {
	Code          uint
	StatusMessage string
	Value         string
}

// authConnection asks a user for his login and password and checks if his own map
// exists in storage, if not then it will be created.
func (s *PotatoSlave) authConnection(connection net.Conn) (string, error) {
	if _, ok := s.storage["user"]; ok {
	} else {
		s.storage["user"] = make(map[string]potat)
	}

	return "user", nil
}

func (s *PotatoSlave) handleConnection(connection net.Conn, username string) {

	defer connection.Close()

	decoder := json.NewDecoder(connection)
	encoder := json.NewEncoder(connection)
	var mes CommandMessage
	for {

		connection.SetReadDeadline(time.Now().Add(s.STALETIME))
		err := decoder.Decode(&mes)

		if err != nil {
			// TODO: check if it's a timeout and then just close the connection
			// Say that you are available
			s.availableWorkers <- true
			return
		}

		returnMes := s.functions[mes.Name](username, mes)
		encoder.Encode(returnMes)

	}
}

//////////
// Invocable functions
//////////

///// Service messages

const (
	_OK = iota
	_WT = iota
	_NK = iota
	_WA = iota
	_NW = iota
)

var statusMessages = map[uint]string{
	_OK: "OK",
	_WT: "Object stored at the key is of different type",
	_NK: "Key doesn't exist",
	_WA: "Wrong call arguments",
	_NW: "There are no available workers on the server",
}

func setStatus(mes *ResponseMessage, code uint) {
	mes.Code = code
	mes.StatusMessage = statusMessages[code]
}

//////////////////////////

///// Data independent Functions

func (s *PotatoSlave) del(userID string, mes CommandMessage) ResponseMessage {

	var response ResponseMessage

	if len(mes.Arguments) != 1 {
		setStatus(&response, _WA)
	} else {
		delete(s.storage[userID], mes.Arguments[0])
		setStatus(&response, _OK)
	}

	return response
}

// TODO: get rid of the boilerplate in here...

//// String functions

func (s *PotatoSlave) get(userID string, mes CommandMessage) ResponseMessage {

	var response ResponseMessage

	if len(mes.Arguments) != 1 {
		setStatus(&response, _WA)
	} else {

		if val, ok := s.storage[userID][mes.Arguments[0]]; ok {

			switch val.(type) {
			case *pstring:
				response.Value, _ = val.getContent("")
				setStatus(&response, _OK)
			default:
				setStatus(&response, _WT)
			}

		} else {
			setStatus(&response, _NK)
		}
	}

	return response
}

func (s *PotatoSlave) set(userID string, mes CommandMessage) ResponseMessage {

	var response ResponseMessage
	var ttl time.Duration

	if len(mes.Arguments) != 2 {
		setStatus(&response, _WA)
	} else {

		delete(s.storage[userID], mes.Arguments[0])

		if mes.TTL != 0 {
			ttl = mes.TTL
		} else {
			ttl = s.DEFAULTTTL
		}

		s.storage[userID][mes.Arguments[0]] = &pstring{
			content:     mes.Arguments[1],
			timeOfDeath: time.Now().Add(ttl),
		}
		setStatus(&response, _OK)
	}

	return response
}

//// List functions

func (s *PotatoSlave) lpush(userID string, mes CommandMessage) ResponseMessage {

	var response ResponseMessage

	if len(mes.Arguments) != 2 {
		// currently we don't support addition of multiple elements...
		setStatus(&response, _WA)
	} else {

		// Key exist and it's of the right type
		if val, ok := s.storage[userID][mes.Arguments[0]]; ok {

			switch val.(type) {
			case *plist:

				s.storage[userID][mes.Arguments[0]].setContent(mes.Arguments[1], "-1")
				setStatus(&response, _OK)
				return response

			default:
			}

		}

		var ttl time.Duration

		if mes.TTL != 0 {
			ttl = mes.TTL
		} else {
			ttl = s.DEFAULTTTL
		}

		s.storage[userID][mes.Arguments[0]] = &plist{
			list:        []string{mes.Arguments[1]},
			timeOfDeath: time.Now().Add(ttl),
		}
		setStatus(&response, _OK)
	}

	return response
}

func (s *PotatoSlave) lset(userID string, mes CommandMessage) ResponseMessage {

	var response ResponseMessage

	if len(mes.Arguments) != 3 {
		setStatus(&response, _WA)
	} else {

		if val, ok := s.storage[userID][mes.Arguments[0]]; ok {

			switch val.(type) {
			case *plist:

				err := s.storage[userID][mes.Arguments[0]].setContent(mes.Arguments[2], mes.Arguments[1])

				if err != nil {
					setStatus(&response, _WA)
				} else {
					setStatus(&response, _OK)
				}

			default:
				setStatus(&response, _WT)
			}

		} else {
			setStatus(&response, _NK)
		}
	}

	return response
}

func (s *PotatoSlave) lget(userID string, mes CommandMessage) ResponseMessage {

	var response ResponseMessage

	if len(mes.Arguments) != 2 {
		setStatus(&response, _WA)
	} else {
		if val, ok := s.storage[userID][mes.Arguments[0]]; ok {

			switch val.(type) {
			case *plist:
				content, err := s.storage[userID][mes.Arguments[0]].getContent(mes.Arguments[1])

				if err != nil {
					setStatus(&response, _WA)
				} else {
					response.Value = content
					setStatus(&response, _OK)
				}

			default:
				setStatus(&response, _WT)
			}
		} else {
			setStatus(&response, _NK)
		}
	}

	return response
}
