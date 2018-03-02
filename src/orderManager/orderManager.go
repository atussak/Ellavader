package orderManager

import(
	"../elevio"
	def "../definitions"
)


type Order struct{
	Button elevio.ButtonType
	Floor int
}

func isQueueEmpty(requests [][]bool) bool {
	for floor := 0; floor < def.NUM_FLOORS; floor++ {
		for button := 0; button < def.NUM_BUTTON_TYPES; button++ {
			if requests[floor][button] {
				return false
			}
		}
	}
	return true
}

func isOrderInFloor(requests [][]bool, floor int) {
	for button := 0; button <= def.NUM_BUTTON_TYPES; button++ {
		if requests[floor][button] {
			return true
		}
	}
	return false
}

func isOrdersAbove(requests [][]bool, current_floor int) bool {
	for floor:= current_floor+1; floor < def.NUM_FLOORS; floor++ {
		if isOrderInFloor(requests, floor) {
			return true
		}
	}
	return false
}

func isOrdersBelow(requests [][]bool, current_floor int) bool {
	for floor:= current_floor-1; floor >= 0; floor-- {
		if isOrderInFloor(requests, floor) {
			return true
		}
	}
	return false
}



func ChooseDirection(requests [][]bool, current_floor int, current_direction elevio.Motor) {

	if isQueueEmpty(requests) { return elevio.MD_Stop }

	if current_direction = elevio.MD_Up {
		if isOrdersAbove(requests, current_floor) {
			return current_direction
		}
	} else if current_direction = elevio.MD_Down {
		if isOrdersBelow(requests, current_floor) {
			return current_direction
		}
	}

	return -1*current_direction
}


