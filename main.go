package main

import (
	"fmt"
	"time"

	"github.com/RXWP5657/DroneProject/Tello-App/tello"
	"github.com/RXWP5657DroneProject/Tello-App/dualshock"
)

func main() {

	drone, err := tello.Init()

	if err != nil {
		fmt.Println("Coulnd`t init drone system: %s", err)
	}

	dualshock.BindController(&drone)

	time.Sleep(60 * time.Second)

	drone.ReleaseDrone()
}
