package main

import (
	"bcast"
	"fmt"
	
)

const localPort = 20011


func main(){
	transmitt := make(chan string)
	receive := make(chan string)
	go bcast.Transmitter(localPort, transmitt)
	go bcast.Receiver(localPort, receive)
	transmitt <- "Hellooo"
	for{
		
	}

}