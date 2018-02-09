package main

import (
	"bcast"
	"fmt"
	
)

const localPort = 20011
const listenPort = 20014

func main(){
	transmitt := make(chan string)
	receive := make(chan string)
	go bcast.Transmitter(localPort, transmitt)
	go bcast.Receiver(listenPort, receive)
	transmitt <- "Hellooo"
	for{
		fmt.Println(<-receive)
	}

}