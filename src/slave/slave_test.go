package slave

import (
	"encoding/json"
	"net"
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

func TestGetSet(t *testing.T) {

	// Create a slave
	testPort := "62554"
	staleTime := time.Second
	s := NewSlave("localhost", testPort, staleTime)

	// Client simulator
	go func(testPort string, s *PotatoSlave, t *testing.T) {

		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", "localhost:"+testPort)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(conn)
		decoder := json.NewDecoder(conn)

		var response ResponseMessage

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
	s.StartServing(true)

}
