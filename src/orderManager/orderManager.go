package orderManager

import(
	"../elevio"
	def "../definitions"
)

var Requests[][] bool 


func IsQueueEmpty() bool {
	for floor := 0; floor < def.NUM_FLOORS; floor++ {
		for button := 0; button < def.NUM_BUTTON_TYPES; button++ {
			if Requests[floor][button] {
				return false
			}
		}
	}
	return true
}

func IsOrderInFloor(floor int) bool{
	for button := 0; button < def.NUM_BUTTON_TYPES; button++ {
		if Requests[floor][button] {
			return true
		}
	}
	return false
}

func IsOrderAbove(current_floor int) bool {
	for floor:= current_floor+1; floor < def.NUM_FLOORS; floor++ {
		if IsOrderInFloor(floor) {
			return true
		}
	}
	return false
}

func IsOrderBelow(current_floor int) bool {
	for floor:= current_floor-1; floor >= 0; floor-- {
		if IsOrderInFloor(floor) {
			return true
		}
	}
	return false
}

func ChooseDirection(current_floor int, current_direction elevio.MotorDirection) elevio.MotorDirection {

	if IsQueueEmpty() { return elevio.MD_Stop }

	if current_direction == elevio.MD_Up {
		if IsOrderAbove(current_floor) {
			return current_direction
		}
	} else if current_direction == elevio.MD_Down {
		if IsOrderBelow(current_floor) {
			return current_direction
		}
	}

	return -1*current_direction
}

func ShouldStopForOrder(order elevio.Order, direction elevio.MotorDirection, current_floor int) bool{
	
	if current_floor != order.Floor { return false}

	if order.Button == elevio.BT_Cab { return true }

	if direction == elevio.MD_Up && IsOrderAbove(current_floor) {
		return order.Button == elevio.BT_HallUp
	} else if direction == elevio.MD_Down && IsOrderBelow(current_floor) {
		return order.Button == elevio.BT_HallDown
	} else {
		return true
	}
	
}

func ShouldStop(direction elevio.MotorDirection, current_floor int) bool{

	execute_cab := ShouldStopForOrder(elevio.Order{current_floor, elevio.BT_Cab}, direction, current_floor)
	execute_up := ShouldStopForOrder(elevio.Order{current_floor, elevio.BT_HallUp}, direction, current_floor)
	execute_down := ShouldStopForOrder(elevio.Order{current_floor, elevio.BT_HallDown}, direction, current_floor)

	return execute_cab && Requests[current_floor][elevio.BT_Cab] ||
		execute_down && Requests[current_floor][elevio.BT_HallDown] ||
		execute_up && Requests[current_floor][elevio.BT_HallUp]

}

func ClearOrder(current_floor int, direction elevio.MotorDirection){

	elevio.SetDoorOpenLamp(false)
	elevio.SetButtonLamp(elevio.BT_Cab, current_floor, false)
	Requests[current_floor][elevio.BT_Cab] = false

	execute_up := ShouldStopForOrder(elevio.Order{current_floor, elevio.BT_HallUp}, direction, current_floor)
	execute_down := ShouldStopForOrder(elevio.Order{current_floor, elevio.BT_HallDown}, direction, current_floor)


	if execute_down {
		elevio.SetButtonLamp(elevio.BT_HallDown, current_floor, false)
		Requests[current_floor][elevio.BT_HallDown] = false
	}
	if execute_up {
		elevio.SetButtonLamp(elevio.BT_HallUp, current_floor, false)
		Requests[current_floor][elevio.BT_HallUp] = false	
	}
	
}