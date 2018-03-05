package main

import (
	"bcast"
	"fmt"
	"time"
	"../localip"
	"../peers"
)

type ElevatorUpdate struct {
	Requests [][] bool,
	Last_floor int,
	Direction elevio.MotorDirection,
	State int,
	ID string,
}


func Init(peer_update_ch chan peers.PeerUpdate, peer_tx_enable chan bool) {

	var id string

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	go peers.Transmitter(15647, id, peer_tx_enable)
	go peers.Receiver(15647, peer_update_ch)
}






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