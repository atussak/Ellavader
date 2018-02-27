package fsm

import(
	"../elevio"
	//"fmt"
)


const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
)


type Channels struct{
	New_order_ch chan elevio.Order
	Direction_ch chan int
	Floor_reached_ch chan int

	Start_timer_ch chan bool
	Timeout_ch chan bool
}


var state int
var direction elevio.MotorDirection
var requests [4] bool
var floor int


func FSM_init(current_floor int){
	floor = current_floor
	state = IDLE
	requests = [4]bool{false, false, false, false}
}


func FSM_run(ch Channels){
	for{
		select {
		// Go routine in main polling buttons	
		case new_order := <- ch.New_order_ch:
			eventNewOrder(new_order, ch)
		case floor = <- ch.Floor_reached_ch:
			eventFloorReached(ch)
		case <- ch.Timeout_ch:
			eventTimeout()
		}
	}
}


func eventNewOrder(new_order elevio.Order, ch Channels){

	// add new order to queue
	requests[new_order.Floor] = true

	// turn on lamp for button pressed
	elevio.SetButtonLamp(new_order.Button, new_order.Floor, true)

	switch state {

	case IDLE:
		if floor == new_order.Floor {
			ch.Start_timer_ch <- true
			state = DOOR_OPEN
		} else {
			if new_order.Floor < floor {
				elevio.SetMotorDirection(elevio.MD_Down)
			} else {
				elevio.SetMotorDirection(elevio.MD_Up)
			}
			
			state = MOVING
		}

	case MOVING:
		// Do nothing

	case DOOR_OPEN:
		if floor == new_order.Floor {
			ch.Start_timer_ch <- true
		}
	}
}


func eventFloorReached(ch Channels){
	// current state is MOVING

	elevio.SetFloorIndicator(floor)

	if requests[floor] {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		ch.Start_timer_ch <- true
		state = DOOR_OPEN
	}
}


func OM_isQueueEmpty() bool {
	for i := 0; i < 4; i++ {
		if requests[i] {
			return true
		}
	}
	return false
}

func OM_chooseDirection() elevio.MotorDirection {

	var dir elevio.MotorDirection = elevio.MD_Stop 

	if OM_isQueueEmpty() {
		return dir
	}

	if direction == elevio.MD_Down {
		for i:= floor-1; i >= 0; i-- {
			if requests[i] {
				dir = elevio.MD_Down
			}
		}
	} else if direction == elevio.MD_Up {
		for i:= floor+1; i < 4; i++ {
			if requests[i] {
				dir =  elevio.MD_Up
			}
		}
	} else if dir == elevio.MD_Stop {
		dir = -1*direction
	}

	return dir
	

}


func eventTimeout(){

	elevio.SetDoorOpenLamp(false)

	requests[floor] = false

	direction = OM_chooseDirection()

	if OM_isQueueEmpty() {
		state = IDLE
	} else {
		direction = OM_chooseDirection()
		state = MOVING
	}
}