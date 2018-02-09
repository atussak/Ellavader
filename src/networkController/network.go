// Network module (exercise 4)

package networkController

import (
	//"fmt"
	"../localip"
	//"bytes"
	//"net"
	//"encoding/gob"
)

var local_IP string

type Packet struct {
	ID string
	data []byte
}


func network_init(receive <-chan Packet, send chan<- Packet)(<-chan Packet, chan<- Packet){
	local_IP, _ = LocalIP()


	
	return receive, send
}

