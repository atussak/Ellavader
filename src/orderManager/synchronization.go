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
	Hall_lights [][] bool
}


var Elevator_database map[int]ElevatorData

var Local_data ElevatorData


func PrintRequests(requests [][] bool) {
	fmt.Printf("Requests: \n")
	for i := 0; i < def.NUM_BUTTON_TYPES; i++ {
		fmt.Printf("\t\t")
		for j := 0; j < def.NUM_FLOORS; j++ {
			if i == int(elevio.BT_HallUp) && j == def.NUM_FLOORS-1 || i==int(elevio.BT_HallDown) && j==0 {
				fmt.Printf("--- \t")
			} else {
				fmt.Printf("%v\t",requests[j][i])
			}
		}
		fmt.Printf("\n\n")
	}
}

func PrintElevatorDatabase() {
	for ID := range Elevator_database {
		fmt.Printf("ID:\t\t%v\n", ID)
		fmt.Printf("State:\t\t%v\n", Elevator_database[ID].State)
		fmt.Printf("Floor:\t\t%v\n", Elevator_database[ID].Floor)
		fmt.Printf("Direction:\t%v\n", Elevator_database[ID].Direction)
		PrintRequests(Elevator_database[ID].Requests)
	}
}



func Init(local_id int) {

	Elevator_database = make(map[int]ElevatorData)

	Local_data = ElevatorData{
		ID: 		local_id,
		State: 		def.IDLE,
		Floor:		0,
		Direction: 	elevio.MD_Down,
		Requests: 	make([][]bool, def.NUM_FLOORS),
		Hall_lights:make([][]bool, def.NUM_FLOORS),
	}

    for floor := 0; floor < def.NUM_FLOORS; floor++ {
        Local_data.Requests[floor] = make([]bool, def.NUM_BUTTON_TYPES)
        for button := 0; button < def.NUM_BUTTON_TYPES; button++ {
            Local_data.Requests[floor][button] = false;
        }
    }

    for floor := 0; floor < def.NUM_FLOORS; floor++ {
        Local_data.Hall_lights[floor] = make([]bool, def.NUM_BUTTON_TYPES-1)
        for button := elevio.BT_HallUp; button <= elevio.BT_HallDown; button++ {
            Local_data.Hall_lights[floor][button] = false;
        }
    }

    Elevator_database[Local_data.ID] = Local_data
}

func UpdateElevatorDatabase(new_elev_data ElevatorData, elev_update_tx_ch chan ElevatorData) {
	Elevator_database[new_elev_data.ID] = new_elev_data
	//PrintElevatorDatabase()
}

func UpdateLocalRequests(floor int, button elevio.ButtonType, value bool, elev_update_tx_ch chan ElevatorData) {
	t := Elevator_database[Local_data.ID]
	t.Requests[floor][button] = value
	Elevator_database[Local_data.ID] = t
	elev_update_tx_ch <- Elevator_database[Local_data.ID]
}

func UpdateLocalState(state int, elev_update_tx_ch chan ElevatorData) {
	t := Elevator_database[Local_data.ID]
	t.State = state
	Elevator_database[Local_data.ID] = t
	elev_update_tx_ch <- Elevator_database[Local_data.ID]
}

func UpdateLocalFloor(floor int, elev_update_tx_ch chan ElevatorData) {
	t := Elevator_database[Local_data.ID]
	t.Floor = floor
	Elevator_database[Local_data.ID] = t
	elev_update_tx_ch <- Elevator_database[Local_data.ID]
}

func UpdateLocalDirection(dir elevio.MotorDirection, elev_update_tx_ch chan ElevatorData) {
	t := Elevator_database[Local_data.ID]
	t.Direction = dir
	Elevator_database[Local_data.ID] = t
	elev_update_tx_ch <- Elevator_database[Local_data.ID]
}

func RegisterOrder(order elevio.Order, elev_update_tx_ch chan ElevatorData) {
	t := Elevator_database[Local_data.ID]
	t.Hall_lights[order.Floor][order.Button] = true
	Elevator_database[Local_data.ID] = t
	elev_update_tx_ch <- Elevator_database[Local_data.ID]
}

func AcceptRemoteOrder(order elevio.Order) {
	
	time.Sleep(10 * time.Millisecond)

	lights_synced := true

	for _, data := range Elevator_database {
		if !data.Hall_lights[order.Floor][order.Button] {
			lights_synced = false;
		}
	}

	if lights_synced {
		elevio.SetButtonLamp(order.Button, order.Floor, true)
	}
}

