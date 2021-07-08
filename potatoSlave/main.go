package main

import (
	"os"
	"potatoSlave/slave"
	"strconv"
	"time"
)

func main() {

	port := os.Getenv("PORT")
	ip := os.Getenv("IP")
	st, _ := strconv.Atoi(os.Getenv("STALETIME"))
	staletime := time.Second * time.Duration(st)
	ttl, _ := strconv.Atoi(os.Getenv("DEFAULTTTL"))
	defaultttl := time.Second * time.Duration(ttl)

	s := slave.NewSlave(ip, port, staletime, defaultttl)
	s.StartServing()
}
