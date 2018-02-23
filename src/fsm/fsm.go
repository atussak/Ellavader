package fsm

import(
	"./elevio"
	"fmt"
)


const (
	IDLE = 0
	MOVING = 1
	IN_FLOOR = 2
)


type Channels struct{
	new_order_ch chan bool
	direction_ch chan int
	in_floor_ch chan int

	start_timer_ch chan bool
	timeout_ch chan bool

}


var current_state int


func FSM(channels Channels){
	for{
		select{
		case <- channels.new_order_ch:
			eventNewOrder()
		case <- channels.in_floor_ch:
			eventInFloor()
		case <- channels.timeout_ch:
			eventTimeout()
		}
	}
}


func eventNewOrder(){

	// finn ut hvilken ordre

	// legg i requests

	// sett på lys

	switch current_state:

	case IDLE:
		current_state = MOVING

	case MOVING:

	case IN_FLOOR:

}


func eventInFloor(){
	// start timer

	// skru av lys

	// sjekke/oppdatere kø

	// sett state (kjør forbi eller stopp)

	// åpne dør hvis bestilling

	//

}


func eventTimeout(){
	// sett state til idle hvis kø er tom

	// sett state til moving hvis kø ikke er tom

	// skru av door open

	//
}