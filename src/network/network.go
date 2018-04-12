package main

import (
	"bcast"
	"fmt"
	"time"
	"../localip"
	"../peers"
)

onlineElevators := 0

func OfflineElevator() (bool,int) {

}





/*

 x -pass på at heis som ikke er på nettverket ikke blir assignet ordre
 x -assign død heis sine ordre til andre
 x -slett hall orders fra den døde heisen
 x -skriv cab orders til disk hver gang noe endres (i cab orders)
 x -les cab orders fra disk i starten av programmet
   -timeout på om ordre blir tatt

   -fikse git repo

TESTING
-test på fysiske heiser
-test med en heis alene
-test med to heiser alene
-test med en heis som dør og kommer tilbake

*/