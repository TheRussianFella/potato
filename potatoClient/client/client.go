package client

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// CommandMessage is a message that we send to the server
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

// Server is a structure that represents a potatoSlave
type Server struct {
	encoder  *json.Encoder
	decoder  *json.Decoder
	response ResponseMessage
}

// Connect
func (s *Server) Connect(path string) {

	conn, err := net.Dial("tcp", path)
	if err != nil {
		panic(err)
	}
	s.encoder = json.NewEncoder(conn)
	s.decoder = json.NewDecoder(conn)
}

// Get
func (s *Server) Get(key string) string {
	s.encoder.Encode(CommandMessage{
		Name:      "GET",
		Arguments: []string{key},
	})
	s.decoder.Decode(&s.response)
	fmt.Println(s.response.StatusMessage)
	return s.response.Value
}

// Set
func (s *Server) Set(key string, value string, ttl time.Duration) {
	s.encoder.Encode(CommandMessage{
		Name:      "SET",
		Arguments: []string{key, value},
		TTL:       ttl,
	})
	s.decoder.Decode(&s.response)
	fmt.Println(s.response.StatusMessage)
}

// Keys
func (s *Server) Keys() string {
	s.encoder.Encode(CommandMessage{
		Name: "KEYS",
	})
	s.decoder.Decode(&s.response)
	fmt.Println(s.response.StatusMessage)
	return s.response.Value
}
