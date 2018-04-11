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
-kan vi sende flere meldinger på en channel så lenge vi har buffer?
-trenger vi å skrive cab orders til disk?

-pass på at heis som ikke er på nettverket ikke blir assignet ordre
-assign død heis sine ordre til andre
-slett hall orders fra den døde heisen
-skriv cab orders til disk hver gang noe endres (i cab orders)
-les cab orders fra disk i starten av programmet
-timeout på om ordre blir tatt

-fikse git repo

TESTING
-test på fysiske heiser
-test med en heis alene
-test med to heiser alene
-test med en heis som dør og kommer tilbake

*/