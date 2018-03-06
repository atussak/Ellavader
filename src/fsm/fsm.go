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

	peer_update_ch

}


var state int
var direction elevio.MotorDirection
var floor int


func Init(current_floor int){
	floor = current_floor
	state = IDLE
	fmt.Printf("IDLE\n")
	//OM.Requests = [4]bool{false, false, false, false}
	//request = make([]bool, def.NUM_FLOORS, def.NUM_FLOORS )

	OM.Requests = make([][]bool, def.NUM_FLOORS)
    for floor := 0; floor < def.NUM_FLOORS; floor++ {
        OM.Requests[floor] = make([]bool, def.NUM_BUTTON_TYPES)
        for button := 0; button < def.NUM_BUTTON_TYPES; button++ {
            OM.Requests[floor][button] = false;
        }
    }


}


func Run(ch Channels){
	for{

		select {
		// Go routine in main polling buttons	
		case new_order := <- ch.New_order_ch:
			//oppdater elevator_database[local_ID]
			fmt.Printf("New order \n")
			eventNewOrder(new_order, ch)
		case floor = <- ch.Floor_reached_ch:
			//oppdater elevator_database[local_ID]
			fmt.Printf("New floor \n")
			eventFloorReached(ch)
		case <- ch.Timeout_ch:
			//oppdater elevator_database[local_ID]
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


func eventNewOrder(new_order elevio.Order, ch Channels){

	// add new order to queue
	OM.Requests[new_order.Floor][new_order.Button] = true

	// turn on lamp for button pressed
	elevio.SetButtonLamp(new_order.Button, new_order.Floor, true)

	switch state {

	case IDLE:
		if floor == new_order.Floor{
			elevio.SetDoorOpenLamp(true)
			ch.Start_timer_ch <- true
			state = DOOR_OPEN
			fmt.Printf("DOOR_OPEN\n")
		} else {
			if new_order.Floor < floor {
				direction = elevio.MD_Down
			} else {
				direction = elevio.MD_Up
			}
			elevio.SetMotorDirection(direction)
			state = MOVING
			fmt.Printf("MOVING\n")
		}

	case MOVING:
		// Do nothing

	case DOOR_OPEN:
		if OM.ShouldStopForOrder(new_order, direction, floor) {
			ch.Start_timer_ch <- true
		} 
	}
}


func eventFloorReached(ch Channels){
	// current state is MOVING

	elevio.SetFloorIndicator(floor)

	if OM.ShouldStop(direction, floor) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		ch.Start_timer_ch <- true
		state = DOOR_OPEN
		fmt.Printf("DOOR_OPEN\n")
	}

	// Change direction in top and bottom floor
	if floor == 0 || floor == def.NUM_FLOORS-1{
		direction = -1*direction
	}
	
}


func eventTimeout(){

	OM.ClearOrder(floor, direction)

	if OM.IsQueueEmpty() {
		fmt.Printf("IDLE \n")
		state = IDLE
	} else {
		//new_dir := OM_chooseDirection()
		new_dir := OM.ChooseDirection(floor, direction)
		if new_dir != elevio.MD_Stop{
			direction = new_dir
		}
		elevio.SetMotorDirection(new_dir)
		state = MOVING
		fmt.Printf("MOVING\n")
	}
	
}