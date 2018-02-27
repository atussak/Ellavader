package main

import (
    "./elevio"
    "./fsm"
    //"./orderManager"
    // "fmt"
)


func main(){


    ch := fsm.Channels {
        New_order_ch:       make(chan elevio.Order),
        Direction_ch:       make(chan int),
        Floor_reached_ch:   make(chan int),

        Start_timer_ch:     make(chan bool),
        Timeout_ch:         make(chan bool),
    }

    fsm.FSM_init(0)

    go elevio.PollButtons(ch.New_order_ch)
    go elevio.PollFloorSensor(ch.Floor_reached_ch)


    /*
    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)  
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    
    

    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            
        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)
            
        }
    }    */
}
