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

	out := make(chan tello.Package, 10)
	go d.InputChan(out)

	cmd := func() {
		//Triangle Button
		controller.On(joystick.SquarePress, func(data interface{}) {
			out <- tello.Package{Cmd: tello.Up, Param: "40"}
		})

		controller.On(joystick.XPress, func(data interface{}) {
			out <- tello.Package{Cmd: tello.Down, Param: "40"}
		})

		controller.On(joystick.CirclePress, func(data interface{}) {
			out <- tello.Package{Cmd: tello.RotateC, Param: "45"}
		})

		//Square Button
		controller.On(joystick.TrianglePress, func(data interface{}) {
			out <- tello.Package{Cmd: tello.RotateCC, Param: "45"}
		})

		controller.On(joystick.LeftX, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -11000: //Left
				out <- tello.Package{Cmd: tello.Left, Param: "20"}
				time.Sleep(15)
			case inData > 11000: // Right
				out <- tello.Package{Cmd: tello.Right, Param: "20"}
				time.Sleep(15)
			}
		})

		controller.On(joystick.LeftY, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -11000: //Up
				out <- tello.Package{Cmd: tello.Forward, Param: "20"}
				time.Sleep(15)
			case inData > 11000: //Down
				out <- tello.Package{Cmd: tello.Back, Param: "20"}
				time.Sleep(15)
			}
		})

		controller.On(joystick.RightX, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -32600: //Up
				out <- tello.Package{Cmd: tello.Takeoff, Param: ""}
				time.Sleep(30)
			case inData > 32600: //Down
				out <- tello.Package{Cmd: tello.Land, Param: ""}
			}
		})

		controller.On(joystick.R2, func(data interface{}) {
			inData := data.(int16)
			if inData > 32700 {
				out <- tello.Package{Cmd: tello.Battery, Param: ""}
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
