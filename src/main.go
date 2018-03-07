package main

import (
    "./elevio"
    "./fsm"
    def "./definitions"
)


func main(){
	
	elevio.Init("localhost:15657", def.NUM_FLOORS)

    ch := fsm.Channels {
        New_order_ch:       make(chan elevio.Order),
        Direction_ch:       make(chan int),
        Floor_reached_ch:   make(chan int),

        Start_timer_ch:     make(chan bool),
        Timeout_ch:         make(chan bool),
    }

    

    go elevio.PollButtons(ch.New_order_ch)
    go elevio.PollFloorSensor(ch.Floor_reached_ch)

    current_floor := <-ch.Floor_reached_ch

    fsm.Init(current_floor)

    go fsm.DoorTimer(ch.Start_timer_ch, ch.Timeout_ch)
    go fsm.Run(ch)


    for{
    	
    }
}
