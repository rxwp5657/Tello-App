package tello

import (
	"bufio"
	"fmt"
	"net"
)

const (
	ip      = "192.168.10.1"       //Dron IP
	cmdPort = "192.168.10.1:8889"  //Send Command and Receive Response Port
	strPort = "192.168.10.1:11111" //Receive Video Stream Port
)

//Command allows to send a instruction to the drone
type Command int

const (

	//InitCommand sets the drone into dev mode
	InitCommand Command = 0

	//Takeoff tells the drone to take off
	Takeoff Command = 1

	//Land tells the drone to land
	Land Command = 2

	//StreamOn tells the drone to start streaming video
	StreamOn Command = 3

	//StreamOff tells the drone to stop streaming video
	StreamOff Command = 4

	//Emergency tells the drone to stop all functions
	Emergency Command = 5

	//Up tells the drone to go up n units
	Up Command = 6

	//Down tells the drone to go down n units
	Down Command = 7

	//Left tells the drone to go left n units
	Left Command = 8

	//Right tells the drone to go right n units
	Right Command = 9

	//Forward tells the drone to go forward n units
	Forward Command = 10

	//Back tells the drone to go backwards n units
	Back Command = 11

	//RotateC tells the drone to rotate clockwise n degrees
	RotateC Command = 12

	//RotateCC tells the drone to rotate counter clockwise n degrees
	RotateCC Command = 13

	//FlipL tells the drone to flip to the left
	FlipL Command = 14

	//FlipR tells the drone to flip to the right
	FlipR Command = 15

	//FlipF tells the drone to flip forwards
	FlipF Command = 16

	//FlipB tells the drone to flip backwards
	FlipB Command = 17

	//Speed command gets the drone's speed
	Speed Command = 19

	//Battery commmand gets the drone battery level
	Battery Command = 20

	//Time command gets the drone flying time
	Time Command = 21

	//WiFi command gets the drone wifi signal
	WiFi Command = 22

	//Height command gets the drone height
	Height Command = 23

	//Temp command gets the temperature of the flying area
	Temp Command = 24

	//Attitude command gets the yawn, pitch and roll data
	Attitude Command = 25

	//Baro command gets barometric data from the drone
	Baro Command = 26

	//Acceleration command gets the acceleration data from the drone
	Acceleration Command = 27

	//TOF command gets distance value from TOF
	TOF Command = 28
)

//Drone structure that represents the drone
type Drone struct {
	cmdConn net.Conn
	strConn net.Conn
}

// Init the command, state and video streaming systems
func Init() (Drone, error) {

	telloDrone := Drone{}

	cmdConn, cmdErr := net.Dial("udp", cmdPort)
	if cmdErr != nil {
		return telloDrone, fmt.Errorf("drone init: unable to stablish connection to port: %s, because: %s", cmdPort, cmdErr)
	}

	telloDrone.cmdConn = cmdConn

	sendCommand("command", cmdConn)

	strConn, strErr := net.Dial("udp", strPort)
	if strErr != nil {
		return telloDrone, fmt.Errorf("drone init: unable to stablish connection to port : %s, because: %s", strPort, strErr)
	}

	telloDrone.strConn = strConn

	return telloDrone, nil
}

// ReleaseDrone command stops all connections to the drone
func (d *Drone) ReleaseDrone() {
	d.cmdConn.Close()
	d.strConn.Close()
}

//doCommand sends the message 'command' with the parameters 'param' to the drone d
func (d Drone) doCommand(cmd Command, param string) {
	switch cmd {
	case InitCommand:
		sendCommand("command", d.cmdConn)
	case Takeoff:
		sendCommand("takeoff", d.cmdConn)
	case Land:
		sendCommand("land", d.cmdConn)
	case StreamOn:
		sendCommand("streamon", d.cmdConn)
	case StreamOff:
		sendCommand("streamoff", d.cmdConn)
	case Emergency:
		sendCommand("emergency", d.cmdConn)
	case Up:
		sendCommand("up "+param, d.cmdConn)
	case Down:
		sendCommand("down "+param, d.cmdConn)
	case Left:
		sendCommand("left "+param, d.cmdConn)
	case Right:
		sendCommand("rigth "+param, d.cmdConn)
	case Forward:
		sendCommand("forward "+param, d.cmdConn)
	case Back:
		sendCommand("back "+param, d.cmdConn)
	case RotateC:
		sendCommand("cw "+param, d.cmdConn)
	case RotateCC:
		sendCommand("ccw "+param, d.cmdConn)
	case FlipL:
		sendCommand("flip l", d.cmdConn)
	case FlipR:
		sendCommand("flip r", d.cmdConn)
	case FlipF:
		sendCommand("flip f", d.cmdConn)
	case FlipB:
		sendCommand("flip b", d.cmdConn)
	case Speed:
		sendCommand("speed?", d.cmdConn)
	case Battery:
		sendCommand("battery", d.cmdConn)
	case Time:
		sendCommand("time?", d.cmdConn)
	case WiFi:
		sendCommand("wifi?", d.cmdConn)
	case Height:
		sendCommand("height?", d.cmdConn)
	case Temp:
		sendCommand("temp?", d.cmdConn)
	case Attitude:
		sendCommand("attitude?", d.cmdConn)
	case Baro:
		sendCommand("baro?", d.cmdConn)
	case Acceleration:
		sendCommand("acceleration?", d.cmdConn)
	case TOF:
		sendCommand("tof?", d.cmdConn)
	default:

	}
}

func sendCommand(cmd string, cn net.Conn) (string, error) {
	_, err := fmt.Fprintf(cn, cmd)
	if err != nil {
		return "error", fmt.Errorf("couldn't send command because %s", err)
	}
	go HandleResponse(cn)
	return cmd, nil
}

// HandleResponse prints to the console the result of the command (ok, err [drone], err [reader])
func HandleResponse(c net.Conn) {
	p := make([]byte, 2048)
	_, err := bufio.NewReader(c).Read(p)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(string(p))
}

func (d Drone) CommandQueue(in chan Command, param string) {
	for {
		cmd := <-in
		d.doCommand(cmd, param)
	}
}
