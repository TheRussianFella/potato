package slave

import (
	"encoding/json"
	"fmt"
	"net"
)

type CommandMessage struct {
	Name      string
	Arguments []string
}

// StartServing begins an infinite loop for serving connections.
func (s *PotatoSlave) StartServing() {

	listener, err := net.Listen("tcp4", ":"+s.port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {

		c, err := listener.Accept()
		if err != nil {
			panic(err)
		} else {

			decoder := json.NewDecoder(c)

			for i := 0; i < 2; i++ {
				var mes CommandMessage
				err := decoder.Decode(&mes)
				if err != nil {
					panic(err)
				}
				s.v = mes.Name
				fmt.Println(mes.Name)
			}

			c.Close()
			return
		}
	}

}
