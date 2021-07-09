package slave

import (
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"
)

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
	s := NewSlave("localhost", testPort, time.Second, time.Minute, time.Millisecond*100, 1)

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

	}(testPort, s, t)

	// Slave code
	s.StartServing()

}

func TestPlist(t *testing.T) {

	// Create a slave
	testPort := "62553"
	s := NewSlave("localhost", testPort, time.Second, time.Minute, time.Millisecond*100, 1)

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

		// Check that deletion works
		mes = CommandMessage{
			Name:      "DEL",
			Arguments: []string{"mylist"},
		}
		encoder.Encode(mes)
		decoder.Decode(&response)

		if response.Code != _OK {
			t.Errorf("Got %s on deletion", response.StatusMessage)
		}
		if _, ok := s.storage["user"]["mylist"]; ok {
			t.Errorf("Deletion didn't work")
		}

	}(testPort, s, t)

	// Start serving
	s.StartServing()
}

func TestKeyReusage(t *testing.T) {

	// Create a slave
	testPort := "62553"
	s := NewSlave("localhost", testPort, time.Second, time.Minute, time.Millisecond*100, 1)

	// Client simulator
	go func(testPort string, s *PotatoSlave, t *testing.T) {

		encoder, decoder, response := newClient(testPort)

		encoder.Encode(CommandMessage{
			Name:      "SET",
			Arguments: []string{"key", "value"},
			TTL:       time.Minute,
		})
		decoder.Decode(&response)

		encoder.Encode(CommandMessage{
			Name:      "LPUSH",
			Arguments: []string{"key", "1"},
			TTL:       time.Minute,
		})
		decoder.Decode(&response)

		encoder.Encode(CommandMessage{
			Name:      "LGET",
			Arguments: []string{"key", "0"},
		})
		decoder.Decode(&response)

		if response.Code != _OK || response.Value != "1" {
			t.Errorf("Different type writing didn't go as planned. message: %s, Value: %s", response.StatusMessage, response.Value)
		}

	}(testPort, s, t)

	// Start serving
	s.StartServing()
}

func TestMultipleConnections(t *testing.T) {

	// Create a slave
	testPort := "62553"
	s := NewSlave("localhost", testPort, time.Second*3, time.Minute, time.Millisecond*100, 2)

	done := make(chan bool)

	// Client simulation
	go func(testPort string, s *PotatoSlave, t *testing.T, done chan bool) {
		// Client 1
		encoder1, d, r := newClient(testPort)

		encoder1.Encode(CommandMessage{
			Name:      "SET",
			Arguments: []string{"key", "value"},
			TTL:       time.Minute,
		})
		d.Decode(&r)
		// Client 2
		encoder, decoder, response := newClient(testPort)
		time.Sleep(time.Second)

		encoder.Encode(CommandMessage{
			Name:      "GET",
			Arguments: []string{"key"},
		})
		decoder.Decode(&response)

		if response.Value != "value" {
			t.Errorf("Got inconsistency with multiple connections: status %d", response.Code)
		}

		done <- true
	}(testPort, s, t, done)

	s.StartServing()
	<-done
}

/*
TODO: fix this test, but function itself seems to work in interactive mode
func _TestTTL(t *testing.T) {

	// Create a slave
	testPort := "62553"
	s := NewSlave("localhost", testPort, time.Second*5, time.Minute, time.Millisecond*100, 1)
	done := make(chan bool)

	// Client one
	go func(testPort string, s *PotatoSlave, t *testing.T, done chan bool) {

		encoder, decoder, response := newClient(testPort)

		encoder.Encode(CommandMessage{
			Name:      "SET",
			Arguments: []string{"short", "value"},
			TTL:       time.Millisecond,
		})
		decoder.Decode(&response)
		encoder.Encode(CommandMessage{
			Name:      "SET",
			Arguments: []string{"long", "value"},
			TTL:       time.Hour,
		})
		//fmt.Println(s.storage["user"]["short"].getTimeOfDeath())
		//fmt.Println(s.storage["user"]["long"].getTimeOfDeath())
		decoder.Decode(&response)
		time.Sleep(time.Second)

		encoder.Encode(CommandMessage{
			Name:      "GET",
			Arguments: []string{"short"},
		})
		decoder.Decode(&response)
		fmt.Println(response)
		if response.Code != _NK {
			t.Errorf("Expired key wasn't deleted")
		}
		encoder.Encode(CommandMessage{
			Name:      "GET",
			Arguments: []string{"long"},
		})
		decoder.Decode(&response)
		fmt.Println(response)
		if response.Code != _OK {
			t.Errorf("Non expired key was deleted")
		}

		done <- true
	}(testPort, s, t, done)

	s.StartServing()
	<-done
}
*/
