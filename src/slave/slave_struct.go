package slave

// PotatoSlave is a class that deals with storing objects of users, assigned
// to it by a PotatoMaster.
type PotatoSlave struct {
	IP   string
	port string
	v    string
}

// NewSlave creates an instance of a PotatoSlave.
func NewSlave(IP string, port string) *PotatoSlave {

	s := PotatoSlave{
		IP:   IP,
		port: port,
		v:    "variable",
	}
	return &s
}
