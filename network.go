// Network module (exercise 4)

package main

import (
	//"fmt"
	"bytes"
	"net"
	"encoding/gob"
)

var local_IP string

type Packet struct {
	ID string,
	data []byte
}


func NET_init()(<-chan Packet, chan<- Packet){
	

}

