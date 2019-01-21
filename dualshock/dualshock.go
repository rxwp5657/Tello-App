package dualshock

import (
	"time"

	"github.com/RXWP5657/DroneProject/tello"
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
	cmd := func() {

		//Triangle Button
		controller.On(joystick.SquarePress, func(data interface{}) {
			d.Up("40")
		})

		controller.On(joystick.XPress, func(data interface{}) {
			d.Down("40")
		})

		controller.On(joystick.CirclePress, func(data interface{}) {
			d.RotateC("45")
		})

		//Square Button
		controller.On(joystick.TrianglePress, func(data interface{}) {
			d.RotateCC("45")
		})

		controller.On(joystick.LeftX, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -11000: //Left
				d.Left("50")
				time.Sleep(15)
			case inData > 11000: // Right
				d.Right("50")
				time.Sleep(15)
			}
		})

		controller.On(joystick.LeftY, func(data interface{}) {
			switch inData := data.(int16); {

			case inData < -11000: //Up
				d.Forward("50")
				time.Sleep(15)
			case inData > 11000: //Down
				d.Back("50")
				time.Sleep(15)
			}
		})

		controller.On(joystick.RightX, func(data interface{}) {

			switch inData := data.(int16); {
			case inData < -32600: //Up
				d.TakekOff()
				time.Sleep(2 * time.Second)
			case inData > 32600: //Down
				d.Land()
			}
		})

		controller.On(joystick.R2, func(data interface{}) {
			inData := data.(int16)
			if inData > 32700 {
				d.GetBattery()
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
