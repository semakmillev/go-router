package main

import (
	"time"
	"net"
	"log"
)

func main (){
	timeout := time.Duration(3 * time.Second)
	_, err := net.DialTimeout("tcp", "localhost:5002", timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
	}
}
