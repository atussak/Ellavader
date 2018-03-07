package orderManager

import(
	"../elevio"
	def "../definitions"
	"fmt"
)



type ElevatorData struct {
	ID 			int
	State 		int
	Floor		int
	Direction	elevio.MotorDirection
	Requests 	[][] bool
}


var elevator_database map[int]ElevatorData

local_data := ElevatorData{
	ID: 		20015,
	State: 		IDLE,
	Floor:		0,
	Direction: 	elevio.MD_Up,
	Requests: 	Requests,
}


func UpdateElevatorDatabase(new_elev_data ElevatorData) {
	fmt.Print("Received elevator update\n")
	remote_order_update, floor, button, value := RemoteOrderUpdate(elevator_database[new_elev_data.ID], new_elev_data)
	if remote_order_update {
		UpdateLocalRequests(floor, button, value)
	}

	elevator_database[new_elev_data.ID] = new_elev_data

	
}

func RemoteOrderUpdate(prev_data ElevatorData, new_data ElevatorData) bool, int, elevio.ButtonType, bool {
	prev_requests := prev_data.Requests
	new_requests := new_data.Requests
	for floor := 0; floor < def.NUM_FLOORS; floor++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_HallDown; button++ {
			if prev_requests[floor][button] != new_requests[floor][button] {
				return true,floor,button,new_requests[floor][button];
			}
		} 
	}
	return false,0,0,false
}

func UpdateLocalRequests(floor int, button elevio.ButtonType, value bool, elev_update_tx_ch chan ElevatorData) {
	elevator_database[local_ID].Requests[floor][button] = value

	elev_update_tx_ch <- elevator_database[local_ID]
}

