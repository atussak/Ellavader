package orderManager

import(
	"../elevio"
	def "../definitions"
)


type Order struct{
	Button elevio.ButtonType
	Floor int
}

func isQueueEmpty(requests []bool) bool {
	for i := 0; i < def.NUM_FLOORS; i++ {
		if requests[i] {
			return false
		}
	}
	return true
}


func isOrdersAbove(requests []bool, current_floor int) bool {
	for i:= floor-1; i >= 0; i-- {
		if requests[i] {
			dir = elevio.MD_Down
		}
	}
}

func isOrdersBelow(requests []bool, current_floor int) bool {
	for i:= floor-1; i >= 0; i-- {
		if requests[i] {
			dir = elevio.MD_Down
		}
	}
}

func ChooseDirection(requests []bool, current_floor int, current_direction elevio.Motor) {

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


