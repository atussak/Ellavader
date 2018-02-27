package fsm

import(
	"./elevio"
	"fmt"
)


const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
)


type Channels struct{
	new_order_ch chan bool
	direction_ch chan int
	in_floor_ch chan int

	start_timer_ch chan bool
	timeout_ch chan bool

}


var current_state int
var requests [4] bool


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
		// if etasje = bestillt etasje
			// slett ordre
			// start timer for åpen dør
			// åpne dør
		current_state = DOOR_OPEN
		// else
		current_state = MOVING
		// sett riktig kjøreretning

	case MOVING:
		// Do nothing

	case DOOR_OPEN:
		// if etasje = bestillt etasje
			// slett ordre
			// start timer for åpen dør
		current_state = IDLE


}


func eventInFloor(){
	// current state is MOVING

	// start timer

	// skru av lys

	// sjekke kø

	// åpne dør hvis bestilling
	// (sett på timer og lys)

	// bytt state til DOOR_OPEN

}


func eventTimeout(){
	// skru av door light

	// oppdatere kø

	// sett state til idle hvis kø er tom

	// sett state til moving hvis kø ikke er tom

}