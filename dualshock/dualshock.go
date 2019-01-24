package dualshock

import (
	"time"

	"github.com/RXWP5657/DroneProject/Tello-App/tello"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

//Controller binded to the drone d
type Controller struct {
	Drone *tello.Drone
	Robot *gobot.Robot
}

//BindController bind some controller to the drone d
func BindController(d *tello.Drone) Controller {
	adapter := joystick.NewAdaptor()
	controller := joystick.NewDriver(adapter, "dualshock4")

	in := make(chan tello.Command)
	param := make(chan string)
	go d.InputChan(in, param)

	cmd := func() {
		//Triangle Button
		controller.On(joystick.SquarePress, func(data interface{}) {
			in <- tello.Up
			param <- "40"
		})

		controller.On(joystick.XPress, func(data interface{}) {
			in <- tello.Down
			param <- "40"
		})

		controller.On(joystick.CirclePress, func(data interface{}) {
			in <- tello.RotateC
			param <- "45"
		})

		//Square Button
		controller.On(joystick.TrianglePress, func(data interface{}) {
			in <- tello.RotateCC
			param <- "45"
		})

		controller.On(joystick.LeftX, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -11000: //Left
				in <- tello.Left
				param <- "1"
				time.Sleep(15)
			case inData > 11000: // Right
				in <- tello.Right
				param <- "1"
				time.Sleep(15)
			}
		})

		controller.On(joystick.LeftY, func(data interface{}) {
			switch inData := data.(int16); {

			case inData < -11000: //Up
				in <- tello.Up
				param <- "1"
				time.Sleep(15)
			case inData > 11000: //Down
				in <- tello.Down
				param <- "1"
				time.Sleep(15)
			}
		})

		controller.On(joystick.RightX, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -32600: //Up
				in <- tello.Takeoff
				param <- ""
				time.Sleep(2 * time.Second)
			case inData > 32600: //Down
				in <- tello.Land
				param <- "1"
			}
		})

		controller.On(joystick.R2, func(data interface{}) {
			inData := data.(int16)
			if inData > 32700 {
				in <- tello.Battery
				param <- ""
				time.Sleep(15)
			}
		})
	}
	rob := gobot.NewRobot("tello controller",
		[]gobot.Connection{adapter},
		[]gobot.Device{controller},
		cmd)

	go rob.Start()
	return Controller{Drone: d, Robot: rob}
}
