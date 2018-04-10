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
        New_order_ch:                   make(chan elevio.Order),
        Direction_ch:                   make(chan int),
        Floor_reached_ch:               make(chan int),

        Start_timer_ch:                 make(chan bool),
        Timeout_ch:                     make(chan bool),
        Elev_update_tx_ch:              make(chan OM.ElevatorData, 100),
        Remote_order_executed_tx_ch:    make(chan elevio.Order),
    }


    peer_tx_enable              := make(chan bool)
    peer_update_ch              := make(chan peers.PeerUpdate) 

    elev_update_rx_ch           := make(chan OM.ElevatorData, 100)

    new_remote_order_tx_ch      := make(chan elevio.Order)
    new_remote_order_rx_ch      := make(chan elevio.Order, 100)
    remote_order_executed_rx_ch := make(chan elevio.Order, 100)

    // Inits

    
    fsm.Init(0, ch)

    // Variables

    peer_port := 15647
    assign_port := 15648
    order_ex_port := 15649
    

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

    go bcast.Transmitter(order_ex_port, ch.Remote_order_executed_tx_ch)
    go bcast.Receiver(order_ex_port, remote_order_executed_rx_ch)


    for{
        select{
        case elev_update := <- elev_update_rx_ch:
            OM.UpdateElevatorDatabase(elev_update, ch.Elev_update_tx_ch)
        
        case new_remote_order := <- new_remote_order_rx_ch:
            ID := OM.AssignOrderToElevator(new_remote_order)
            // synchronize lights
            OM.AcceptRemoteOrder(new_remote_order)
            if OM.Local_data.ID == ID {
                ch.New_order_ch <- new_remote_order
            }
        case executed_order := <- remote_order_executed_rx_ch:
            elevio.SetButtonLamp(executed_order.Button, executed_order.Floor, false)
        }   
    }
}