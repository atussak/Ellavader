package main

import (
	"bcast"
	"fmt"
	"time"
)

const localPort = 20007


func main(){
	transmit := make(chan string)
	receive := make(chan string)
	go bcast.Transmitter(localPort, transmit)
	go bcast.Receiver(localPort, receive)
	go func(){
		for {
			transmit <- "Hellooo"
			time.Sleep(time.Second)
		}
	}()

	for{
		select {
		case msg := <-receive:
			fmt.Println(msg)

		}
	}
}