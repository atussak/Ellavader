package fsm

import(
	"../elevio"
	"time"
	"fmt"
	def "../definitions"
	OM "../orderManager"
)



type Channels struct{
	New_order_ch chan elevio.Order
	Direction_ch chan int
	Floor_reached_ch chan int

	Start_timer_ch chan bool
	Timeout_ch chan bool
	Elev_update_tx_ch chan OM.ElevatorData
}


var floor int


func Init(current_floor int, ch Channels){
	OM.UpdateLocalFloor(current_floor, ch.Elev_update_tx_ch)
	OM.UpdateLocalState(def.IDLE, ch.Elev_update_tx_ch)

	fmt.Printf("IDLE\n")
}


func Run(ch Channels){
	for{

		select {
		// Go routine in main polling buttons	
		case new_order := <- ch.New_order_ch:
			fmt.Printf("New order \n")
			eventNewOrder(new_order, ch)
		case floor := <- ch.Floor_reached_ch:
			OM.UpdateLocalFloor(floor, ch.Elev_update_tx_ch)
			fmt.Printf("New floor \n")
			eventFloorReached(ch)
		case <- ch.Timeout_ch:
			fmt.Printf("Timeout \n")
			eventTimeout(ch)
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
	OM.Elevator_database[def.LOCAL_ID].Requests[new_order.Floor][new_order.Button] = true
	OM.UpdateLocalRequests(new_order.Floor, new_order.Button, true, ch.Elev_update_tx_ch)
	// turn on lamp for button pressed
	elevio.SetButtonLamp(new_order.Button, new_order.Floor, true)

	switch OM.Elevator_database[def.LOCAL_ID].State {

	case def.IDLE:
		if OM.Elevator_database[def.LOCAL_ID].Floor == new_order.Floor{
			elevio.SetDoorOpenLamp(true)
			ch.Start_timer_ch <- true
			OM.UpdateLocalState(def.DOOR_OPEN, ch.Elev_update_tx_ch)
			fmt.Printf("DOOR_OPEN\n")
		} else {
			if new_order.Floor < OM.Elevator_database[def.LOCAL_ID].Floor {
				OM.UpdateLocalDirection(elevio.MD_Down, ch.Elev_update_tx_ch)
			} else {
				OM.UpdateLocalDirection(elevio.MD_Up, ch.Elev_update_tx_ch)
			}
			elevio.SetMotorDirection(OM.Elevator_database[def.LOCAL_ID].Direction)
			OM.UpdateLocalState(def.MOVING, ch.Elev_update_tx_ch)
			fmt.Printf("MOVING\n")
		}

	case def.MOVING:
		// Do nothing

	case def.DOOR_OPEN:
		if OM.ShouldStopForOrder(new_order, OM.Elevator_database[def.LOCAL_ID].Direction, OM.Elevator_database[def.LOCAL_ID].Floor) {
			ch.Start_timer_ch <- true
		} 
	}
}


func eventFloorReached(ch Channels){
	// current state is MOVING

	elevio.SetFloorIndicator(OM.Elevator_database[def.LOCAL_ID].Floor)

	if OM.ShouldStop(OM.Elevator_database[def.LOCAL_ID].Direction, OM.Elevator_database[def.LOCAL_ID].Floor) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		ch.Start_timer_ch <- true
		OM.UpdateLocalState(def.DOOR_OPEN, ch.Elev_update_tx_ch)
		fmt.Printf("DOOR_OPEN\n")
	}

	// Change direction in top and bottom floor
	if OM.Elevator_database[def.LOCAL_ID].Floor == 0 || OM.Elevator_database[def.LOCAL_ID].Floor == def.NUM_FLOORS-1{
		OM.UpdateLocalDirection(-1*OM.Elevator_database[def.LOCAL_ID].Direction, ch.Elev_update_tx_ch)
	}
	
}


func eventTimeout(ch Channels){

	OM.ClearOrder(OM.Elevator_database[def.LOCAL_ID].Floor, OM.Elevator_database[def.LOCAL_ID].Direction, ch.Elev_update_tx_ch)

	if OM.IsQueueEmpty() {
		fmt.Printf("IDLE \n")
		OM.UpdateLocalState(def.IDLE, ch.Elev_update_tx_ch)
	} else {
		//new_dir := OM_chooseDirection()
		new_dir := OM.ChooseDirection(OM.Elevator_database[def.LOCAL_ID].Floor, OM.Elevator_database[def.LOCAL_ID].Direction)
		if new_dir != elevio.MD_Stop{
			OM.UpdateLocalDirection(new_dir, ch.Elev_update_tx_ch)
		}
		elevio.SetMotorDirection(new_dir)
		OM.UpdateLocalState(def.MOVING, ch.Elev_update_tx_ch)
		fmt.Printf("MOVING\n")
	}
	
}