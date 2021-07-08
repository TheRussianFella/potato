package slave

import (
	"time"
)

// potat is an interface for objects that could be stored in a PotatoSlave's data storage
type potat interface {
	getTimeOfDeath() time.Time
	getContent(string) string
	setContent(string, string)
}

type pstring struct {
	content     string
	timeOfDeath time.Time
}

func (p pstring) getTimeOfDeath() time.Time {
	return p.timeOfDeath
}

func (p pstring) getContent(idx string) string {
	return p.content
}

func (p pstring) setContent(val string, idx string) {}

// PotatoSlave is a class that deals with storing objects of users, assigned
// to it by a PotatoMaster.
type PotatoSlave struct {

	// Constants
	IP   string
	port string

	STALETIME time.Duration

	// Functions - a map that holds invocable functions
	functions map[string]func(string, CommandMessage) ResponseMessage

	// Data - the structure is a nested map, where first level is a separation by users
	// (each user's keys are stored in a separate table) and then a data map itself.
	storage map[string]map[string]potat

	// This is a flag that tells a slave to stop accepting new connections
	// TODO: Find a better solution to this? Is it safe?
	stop bool
}

// NewSlave creates an instance of a PotatoSlave.
func NewSlave(IP string, port string, STALETIME time.Duration) *PotatoSlave {

	s := PotatoSlave{
		IP:        IP,
		port:      port,
		STALETIME: STALETIME,
		storage:   make(map[string]map[string]potat),
		functions: make(map[string]func(string, CommandMessage) ResponseMessage),
		stop:      false,
	}

	s.functions["GET"] = s.get
	s.functions["SET"] = s.set

	return &s
}
