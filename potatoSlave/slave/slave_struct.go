package slave

import (
	"errors"
	"strconv"
	"time"
)

/////////
// Main structure
/////////

// PotatoSlave is a class that deals with storing objects of users, assigned
// to it by a PotatoMaster.
type PotatoSlave struct {

	// Constants
	IP   string
	port string

	STALETIME  time.Duration
	DEFAULTTTL time.Duration
	NUMWORKERS int

	// Functions - a map that holds invocable functions
	functions map[string]func(string, CommandMessage) ResponseMessage

	// Data - the structure is a nested map, where first level is a separation by users
	// (each user's keys are stored in a separate table) and then a data map itself.
	storage map[string]map[string]potat

	// TODO: This is maximum number of connections that server is allowed to open -
	// it's just a hack so that we can easily stop the server for the tests
	numToServ int
	// availableWorkers is a channel that holds NUMWORKERS 1's when new worker is
	// started it takes one element with him and puts it back when it finishes.
	// It's a semaphore, basicly.
	availableWorkers chan bool
}

// NewSlave creates an instance of a PotatoSlave.
func NewSlave(IP string, port string, STALETIME time.Duration, DEFAULTTTL time.Duration, numToServ int) *PotatoSlave {

	nw := 5
	s := PotatoSlave{
		IP:               IP,
		port:             port,
		STALETIME:        STALETIME,
		DEFAULTTTL:       DEFAULTTTL,
		NUMWORKERS:       nw,
		storage:          make(map[string]map[string]potat),
		functions:        make(map[string]func(string, CommandMessage) ResponseMessage),
		numToServ:        numToServ,
		availableWorkers: make(chan bool, nw),
	}

	s.functions["GET"] = s.get
	s.functions["SET"] = s.set
	s.functions["LGET"] = s.lget
	s.functions["LSET"] = s.lset
	s.functions["LPUSH"] = s.lpush
	s.functions["DEL"] = s.del

	for i := 0; i < s.NUMWORKERS; i++ {
		s.availableWorkers <- true
	}

	return &s
}

/////////
// Structures that represent data
/////////

// potat is an interface for objects that could be stored in a PotatoSlave's data storage
type potat interface {
	getTimeOfDeath() time.Time
	getContent(string) (string, error)
	setContent(string, string) error
}

///// String

type pstring struct {
	content     string
	timeOfDeath time.Time
}

func (p *pstring) getTimeOfDeath() time.Time {
	return p.timeOfDeath
}

func (p *pstring) getContent(idx string) (string, error) {
	return p.content, nil
}

func (p *pstring) setContent(val string, idx string) error { return nil }

///// List

type plist struct {
	list        []string
	timeOfDeath time.Time
}

func (p *plist) getTimeOfDeath() time.Time {
	return p.timeOfDeath
}

// getContent retrieves a string on idx position of list.
func (p *plist) getContent(idx string) (string, error) {

	i, err := strconv.Atoi(idx)
	if err != nil {
		return "", errors.New("wr")
	}

	if i >= len(p.list) {
		return "", errors.New("ou")
	}

	return p.list[i], nil
}

// setContent does quite a self explainatory thing. -1 here means a tale of the list.
func (p *plist) setContent(val string, idx string) error {

	i, err := strconv.Atoi(idx)
	if err != nil {
		return errors.New("wr")
	}

	if i == -1 {
		p.list = append(p.list, val)
		return nil
	}

	if i >= len(p.list) {
		return errors.New("ou")
	}

	p.list[i] = val
	return nil

}
