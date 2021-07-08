package client

import (
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
