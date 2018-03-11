package main

import (
    "./elevio"
    "./fsm"
    def "./definitions"
    OM "./orderManager"
    "./peers"
    "./bcast"
    //"./localip"
    "fmt"
    "os"
    "strconv"
)


func main(){

    ip := os.Args[1]
	host := fmt.Sprintf("localhost:%s",ip)
    local_port,_ := strconv.Atoi(ip)


    elevio.Init(host, def.NUM_FLOORS)
    OM.Init(local_port)

    // Channels

    ch := fsm.Channels {
        New_order_ch:       make(chan elevio.Order),
        Direction_ch:       make(chan int),
        Floor_reached_ch:   make(chan int),

        Start_timer_ch:     make(chan bool),
        Timeout_ch:         make(chan bool),
        Elev_update_tx_ch:  make(chan OM.ElevatorData, 100),
    }


    peer_tx_enable := make(chan bool)
    peer_update_ch := make(chan peers.PeerUpdate) 

    elev_update_rx_ch := make(chan OM.ElevatorData, 100)

    new_remote_order_tx_ch := make(chan elevio.Order)
    new_remote_order_rx_ch := make(chan elevio.Order, 100)

    // Inits

    
    fsm.Init(0, ch)

    // Variables

    peer_port := 15647
    assign_port := 15648
    

    // Goroutines

    go elevio.PollCabButtons(ch.New_order_ch)
    go elevio.PollHallButtons(new_remote_order_tx_ch)
    go elevio.PollFloorSensor(ch.Floor_reached_ch)


        // Start FSM

    go fsm.DoorTimer(ch.Start_timer_ch, ch.Timeout_ch)
    go fsm.Run(ch)

    
        // Elevator communication

    go peers.Transmitter(peer_port, ip, peer_tx_enable)
    go peers.Receiver(peer_port, peer_update_ch)

    go bcast.Transmitter(peer_port, ch.Elev_update_tx_ch)
    go bcast.Receiver(peer_port, elev_update_rx_ch)

    go bcast.Transmitter(assign_port, new_remote_order_tx_ch)
    go bcast.Receiver(assign_port, new_remote_order_rx_ch)


    for{
        select{
        case elev_update := <- elev_update_rx_ch:
            OM.UpdateElevatorDatabase(elev_update, ch.Elev_update_tx_ch)
        
        case new_remote_order := <- new_remote_order_rx_ch:
            ID := OM.AssignOrderToElevator(new_remote_order)
            if OM.Local_data.ID == ID {
                ch.New_order_ch <- new_remote_order
            }
        }   
    }
}
