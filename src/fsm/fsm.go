package fsm

import(
	"../elevio"
	"time"
	"fmt"
	def "../definitions"
	OM "../orderManager"
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
var floor int
//var requests [def.NUM_FLOORS] bool
//var requests [] bool //slice > array :)
var requests[][] bool 


func Init(current_floor int){
	floor = current_floor
	state = IDLE
	//requests = [4]bool{false, false, false, false}
	//request = make([]bool, def.NUM_FLOORS, def.NUM_FLOORS )

	requests = make([][]bool, def.NUM_FLOORS)
    for i := 0; floor < def.NUM_FLOORS; i++ {
        request = make([]bool, def.NUM_BUTTON_TYPES)
        for j := 0; button < def.NUM_BUTTON_TYPES; j++ {
            request[floor][button] = false;
        }
    }


}


func Run(ch Channels){
	for{

		select {
		// Go routine in main polling buttons	
		case new_order := <- ch.New_order_ch:
			fmt.Printf("New order \n")
			eventNewOrder(new_order, ch)
		case floor = <- ch.Floor_reached_ch:
			fmt.Printf("New floor \n")
			eventFloorReached(ch)
		case <- ch.Timeout_ch:
			fmt.Printf("Timeout \n")
			eventTimeout()
		}
	}
}


func DoorTimer( Start_timer_ch chan bool, Timeout_ch chan bool) {
	const door_open_dur = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-Start_timer_ch:
			timer.Reset(door_open_dur)
		case <-timer.C:
			timer.Stop()
			Timeout_ch <- true
		}
	}
}

func OM_isQueueEmpty() bool {
	for i := 0; i < def.NUM_FLOORS; i++ {
		if requests[i] {
			return false
		}
	}
	return true
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
		for i:= floor+1; i < def.NUM_FLOORS; i++ {
			if requests[i] {
				dir =  elevio.MD_Up
			}
		}
	} 
	if dir == elevio.MD_Stop {
		dir = -1*direction
	}

	return dir
	

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
				direction = elevio.MD_Down
			} else {
				direction = elevio.MD_Up
			}
			elevio.SetMotorDirection(direction)
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
	if floor == 0 || floor == def.NUM_FLOORS-1{
		direction = -1*direction
	}
	
}



func eventTimeout(){

	elevio.SetDoorOpenLamp(false)
	elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)

	requests[floor] = false

	//if OM_isQueueEmpty() {
	if OM.isQueueEmpty(requests) {
		fmt.Printf("IDLE \n")
		state = IDLE
	} else {
		//new_dir := OM_chooseDirection()
		new_dir := OM.ChooseDirection(requests, floor, direction)
		if new_dir != elevio.MD_Stop{
			direction = new_dir
		}
		elevio.SetMotorDirection(new_dir)
		state = MOVING
	}
}