package slave

import (
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"
)

/*
TODO: Turn it back on when you make a workpool
func TestStale(t *testing.T) {

	// Create a slave
	testPort := "62554"
	staleTime := time.Second
	s := NewSlave("localhost", testPort, staleTime)

	// Send him a message

	go func(testPort string, s *PotatoSlave, staleTime time.Duration, t *testing.T) {

		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", "localhost:"+testPort)
		if err != nil {
			panic(err)
		}
		time.Sleep(staleTime)
		one := make([]byte, 1)
		_, err = conn.Read(one)

		if err != nil {
			t.Errorf("Connection wasn't timed out")
		}
		s.stop = true
	}(testPort, s, staleTime, t)

	// Check them
	s.StartServing(true)
	fmt.Println("skidadle")
}
*/

func newClient(testPort string) (*json.Encoder, *json.Decoder, ResponseMessage) {

	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "localhost:"+testPort)
	if err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	var response ResponseMessage

	return encoder, decoder, response
}

func TestPstring(t *testing.T) {

	// Create a slave
	testPort := "62554"
	s := NewSlave("localhost", testPort, time.Second, time.Minute)

	// Client simulator
	go func(testPort string, s *PotatoSlave, t *testing.T) {

		encoder, decoder, response := newClient(testPort)

		mes1 := CommandMessage{
			Name:      "SET",
			Arguments: []string{"key", "value"},
			TTL:       time.Minute,
		}

		encoder.Encode(mes1)
		decoder.Decode(&response)

		if response.Code != _OK {
			t.Errorf("Set unsuccessful")
		}

		mes2 := CommandMessage{
			Name:      "GET",
			Arguments: []string{"key"},
		}

		encoder.Encode(mes2)
		decoder.Decode(&response)

		if response.Code != _OK {
			t.Errorf("Get unsuccessful")
		}
		if response.Value != "value" {
			t.Errorf("Got wrong value from get")
		}

		s.stop = true
	}(testPort, s, t)

	// Slave code
	s.StartServing()

}

func TestPlist(t *testing.T) {

	// Create a slave
	testPort := "62553"
	s := NewSlave("localhost", testPort, time.Second, time.Minute)

	// Client simulator
	go func(testPort string, s *PotatoSlave, t *testing.T) {

		encoder, decoder, response := newClient(testPort)

		// Make some insertions
		for i := 0; i < 5; i++ {
			mes := CommandMessage{
				Name:      "LPUSH",
				Arguments: []string{"mylist", strconv.Itoa(i)},
			}

			encoder.Encode(mes)
			decoder.Decode(&response)

			if response.Code != _OK {
				t.Errorf("Couldn't push")
			}

			mes = CommandMessage{
				Name:      "LGET",
				Arguments: []string{"mylist", strconv.Itoa(i)},
			}

			encoder.Encode(mes)
			decoder.Decode(&response)

			if response.Code != _OK {
				t.Errorf("Couldn't get. Got error message: %s", response.StatusMessage)
			}
			if response.Value != strconv.Itoa(i) {
				t.Errorf("Got wrong value")
			}
		}

		// Check that changing at index works
		mes := CommandMessage{
			Name:      "LSET",
			Arguments: []string{"mylist", "3", "10"},
		}
		encoder.Encode(mes)
		decoder.Decode(&response)
		mes = CommandMessage{
			Name:      "LGET",
			Arguments: []string{"mylist", "3"},
		}
		encoder.Encode(mes)
		decoder.Decode(&response)
		if response.Value != "10" {
			t.Errorf("Got wrong value after LSET.")
		}

		s.stop = true
	}(testPort, s, t)

	// Start serving
	s.StartServing()
}
