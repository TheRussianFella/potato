package slave

import (
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {

	// Create a slave
	testPort := "62554"
	s := NewSlave("localhost", testPort)

	// Send him a message

	go func(testPort string) {

		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", "localhost:"+testPort)
		if err != nil {
			panic(err)
		}
		mes1 := CommandMessage{
			Name:      "GET",
			Arguments: []string{"a"},
		}
		mes2 := CommandMessage{
			Name:      "POST",
			Arguments: []string{"a"},
		}

		encoder := json.NewEncoder(conn)
		encoder.Encode(mes1)
		encoder.Encode(mes2)
	}(testPort)

	// Check them
	s.StartServing()
	t.Errorf(s.v)
}
