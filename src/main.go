package main

import (
    "./elevio"
    "./fsm"
    def "./definitions"
    "../localip"
    OM "./orderManager"
    //"fmt"
)


func main(){
	
	elevio.Init("localhost:15657", def.NUM_FLOORS)

    // Variables

    peer_port := 15647 //anything?
    local_port := 20015

    current_floor := <-ch.Floor_reached_ch


    // Channels

    ch := fsm.Channels {
        New_order_ch:       make(chan elevio.Order),
        Direction_ch:       make(chan int),
        Floor_reached_ch:   make(chan int),

        Start_timer_ch:     make(chan bool),
        Timeout_ch:         make(chan bool),
    }


    peer_tx_enable := make(chan bool)
    peer_update_ch := make(chan peers.PeerUpdate) 

    elev_update_tx_ch := make(chan OM.ElevatorData)
    elev_update_rx_ch := make(chan OM.ElevatorData)


    // Goroutines

    go elevio.PollButtons(ch.New_order_ch)
    go elevio.PollFloorSensor(ch.Floor_reached_ch)

    go fsm.DoorTimer(ch.Start_timer_ch, ch.Timeout_ch)
    go fsm.Run(ch)

    go peers.Transmitter(peer_port, peer_tx_enable)
    go peers.Receiver(peer_port, peer_update_ch)
    go bcast.Transmitter(local_port, elev_update_tx_ch)
    go bcast.Receiver(local_port, elev_update_rx_ch)


    for{
        select{
        case elev_update := <- elev_update_rx_ch:
            OM.UpdateElevatorDatabase(elev_update)
        }

    }
}
