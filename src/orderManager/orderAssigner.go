package orderManager

import(
	"../elevio"
	def "../definitions"
	"fmt"
)


func findElevatorInDirection(order elevio.Order, current_floor int) (bool, int){
	ID := def.MAXIMUM_ID
	elevator_found := false

	for id, data := range Elevator_database {
		fmt.Printf("1b:\tChecking elevator %v", id)
		if data.Floor == current_floor && isOrderInDirection(order, data){
			if id < ID {
				ID = id
				elevator_found = true
				fmt.Printf("\tgood")
			}
		}
		fmt.Printf("\n")
	}
	return elevator_found, ID
}

func findIdleElevatorInFloor(current_floor int) (bool, int) {
	ID := def.MAXIMUM_ID
	elevator_found := false

	for id, data := range Elevator_database {
		fmt.Printf("1a:\tChecking elevator %v", id)
		if data.Floor == current_floor && data.State == def.IDLE{
			if id < ID {
				ID = id
				elevator_found = true
				fmt.Printf("\tgood")
			}
		}
		fmt.Printf("\n")
	}
	return elevator_found, ID
}

func countOrders(data ElevatorData) int {
	orders := 0
	for floor := def.BOTTOM_FLOOR; floor <= def.TOP_FLOOR; floor++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
			if data.Requests[floor][button] { orders++ }
		}
	}
	return orders
}

// Takes in an order an chooses the best suited elevator to
// execute the order. The function returns the ID of this
// elevator.
func AssignOrderToElevator(order elevio.Order) int {
	// Default ID is set to higher than any ID so that it is
	// guranteed to be replaced by an actual elevator ID

	ID := def.MAXIMUM_ID
	elevator_found := false


	// Iterate through map and find elevator in closest floor
	// going in the right direction
	// If several elevators: choose lowest ID

	top_reached := false
	bottom_reached := false

	for i := def.BOTTOM_FLOOR; i <= def.TOP_FLOOR; i++ {

		// OBS! -ok- Er ikke vits å finne elevator in direction hvis den likevel må snu for å ta ordren!
		// OBS! -ok- Heisen er ikke i data.Floor hvis den også er moving!
		// OBS! -ok- Det er ikke vits å bry seg om retningen til en heis uten ordre
		//			 Da vil avstand være viktigst			 

		if !top_reached {
			current_floor := order.Floor + i
			top_reached = (current_floor == def.TOP_FLOOR)
			elevator_found, ID = findIdleElevatorInFloor(current_floor)
			if !elevator_found {
				elevator_found, ID = findElevatorInDirection(order, current_floor)
			}
			if elevator_found {
				fmt.Printf("1:\tFound elevator %v in floor %v\n", ID, current_floor)
				return ID
			}
		}

		if !bottom_reached {
			current_floor := order.Floor - i
			bottom_reached = (current_floor == def.BOTTOM_FLOOR)
			elevator_found, ID = findIdleElevatorInFloor(current_floor)
			if !elevator_found {
				elevator_found, ID = findElevatorInDirection(order, current_floor)
			}
			if elevator_found {
				fmt.Printf("1:\tFound elevator %v in floor %v\n", ID, current_floor)
				return ID
			}
		}
	}


	// Iterate through map and find elevator with least amount
	// of orders
	// If several elevators: choose lowest ID

	// Least orders so far starting at the maximum amount of possible orders
	// Subtracted by 2 because there is only one hall button in top and bottom floor
	least_orders := def.NUM_FLOORS*def.NUM_BUTTON_TYPES - 2;
	for id, data := range Elevator_database {
		orders := countOrders(data)
		if orders <= least_orders && id < ID {
			least_orders = orders
			ID = id
		}
	}
	fmt.Printf("2:\tFound elevator %v with %v orders.\n", ID, least_orders)
	return ID
}