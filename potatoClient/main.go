package main

import (
	"encoding/json"
	"net"
	"potatoClient/client"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:65000")
	if err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(conn)

	mes1 := client.CommandMessage{
		Name:      "SET",
		Arguments: []string{"key", "value"},
		TTL:       time.Minute,
	}

	encoder.Encode(mes1)
}
