package slave

import (
	"net"
)

// StartServing begins an infinite loop for serving connections.
func (s *PotatoSlave) StartServing(debug bool) {

	listener, err := net.Listen("tcp4", ":"+s.port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for !s.stop {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		name, _ := s.authConnection(c)
		s.handleConnection(c, name)
	}
}
