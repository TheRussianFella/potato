package main

import (
	"fmt"
	"potatoClient/client"
)

func main() {

	serv := client.Server{}
	serv.Connect("localhost:65000")

	//serv.Set("tasemp", "s", time.Second)
	fmt.Println(serv.Keys())
}
