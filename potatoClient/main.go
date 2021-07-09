package main

import (
	"fmt"
	"potatoClient/client"
	"time"
)

func main() {

	serv := client.Server{}
	serv.Connect("localhost:65000")

	// String
	serv.Set("newstr", "value", time.Hour)
	fmt.Println(serv.Get("newstr"))
	serv.Set("superstr", "value", time.Hour)
	fmt.Println(serv.Keys())
	serv.Del("superstr")
	fmt.Println(serv.Keys())

	// List
	serv.Lpush("mylist", "5", time.Hour)
	serv.Lpush("mylist", "10", time.Hour)
	serv.Lpush("mylist", "15", time.Hour)
	serv.Lset("mylist", 2, "22")
	fmt.Println(serv.Lget("mylist", 0))
	fmt.Println(serv.Lget("mylist", 2))
	serv.Lset("notexistent", 2, "22")
	serv.Lset("mylist", 100, "22")

	// Map
	serv.Hset("myhash", "max", "boy", time.Hour)
	serv.Hset("myhash", "nastya", "girl", time.Hour)
	fmt.Println(serv.Hget("myhash", "max"))
	//TTL
	fmt.Println(serv.Keys())
	serv.Set("shortLiver", "value", time.Millisecond)
	fmt.Println(serv.Keys())
	time.Sleep(time.Second)
	fmt.Println(serv.Keys())
}
