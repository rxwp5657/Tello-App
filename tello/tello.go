package tello

import (
	"bufio"
	"fmt"
	"net"
)

const (
	ip      = "192.168.10.1"      //Dron IP
	cmdPort = "192.168.10.1:8889" //Send Command and Receive Response Port
	strPort = "0.0.0.0:11111"     //Receive Video Stream Port
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

	initCmdMode(cmdConn)

	strConn, strErr := net.Dial("udp", strPort)
	if strErr != nil {
		return telloDrone, fmt.Errorf("drone init: unable to stablish connection to port : %s, because: %s", strPort, strErr)
	}

	telloDrone.strConn = strConn

	return telloDrone, nil
}

// ReleaseDrone commands stops all connections to the drone
func (d *Drone) ReleaseDrone() {
	d.cmdConn.Close()
	d.strConn.Close()
}

func handleResponse(c net.Conn) {
	p := make([]byte, 2048)
	_, err := bufio.NewReader(c).Read(p)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(string(p))
}

func initCmdMode(c net.Conn) (string, error) {
	return sendCommand("command", c)
}

//Emergency command set the drone into Emergency mode
func (d Drone) Emergency() (string, error) {
	return sendCommand("emergency", d.cmdConn)
}

//TakekOff command
func (d Drone) TakekOff() (string, error) {
	return sendCommand("takeoff", d.cmdConn)
}

//Land command
func (d Drone) Land() (string, error) {
	return sendCommand("land", d.cmdConn)
}

//Up command moves the drone up by Xcm (20cm default)
func (d Drone) Up(x string) (string, error) {
	return sendCommand(("up " + x), d.cmdConn)
}

//Down command moves the drone down by Xcm (20cm default)
func (d Drone) Down(x string) (string, error) {
	return sendCommand(("down " + x), d.cmdConn)
}

//Left command moves the drone to the left by Xcm (20cm default)
func (d Drone) Left(x string) (string, error) {
	return sendCommand(("left " + x), d.cmdConn)
}

//Right command moves the drone to the right by Xcm (20cm default)
func (d Drone) Right(x string) (string, error) {
	return sendCommand(("right " + x), d.cmdConn)
}

//Forward command moves the drone forward by Xcm (20cm default)
func (d Drone) Forward(x string) (string, error) {
	return sendCommand(("forward " + x), d.cmdConn)
}

//Back command moves the drone backwards by Xcm (20cm default)
func (d Drone) Back(x string) (string, error) {
	return sendCommand(("back " + x), d.cmdConn)
}

//RotateC command rotates the drone clockwise by N degrees (10 degrees default)
func (d Drone) RotateC(degrees string) (string, error) {
	return sendCommand(("cw " + degrees), d.cmdConn)
}

//RotateCC command rotates the drone counter clockwise by N degrees (10 degrees default)
func (d Drone) RotateCC(degrees string) (string, error) {
	return sendCommand(("ccw " + degrees), d.cmdConn)
}

//FlipL command flips the drone to the left
func (d Drone) FlipL() (string, error) {
	return sendCommand("flip l", d.cmdConn)
}

//FlipR command flips the drone to the right
func (d Drone) FlipR() (string, error) {
	return sendCommand("flip r", d.cmdConn)
}

//FlipF command flips the drone forwards
func (d Drone) FlipF() (string, error) {
	return sendCommand("flip f", d.cmdConn)
}

//FlipB command flips the drone to the left
func (d Drone) FlipB() (string, error) {
	return sendCommand("flip b", d.cmdConn)
}

//GoXYZ command sends the drone to the XYZ coordinate at some speed S
func (d Drone) GoXYZ(x, y, z, speed string) (string, error) {
	return sendCommand((x + " " + y + " " + z + " " + speed), d.cmdConn)
}

//Speed command sets the drone's speed
func (d Drone) Speed(speed string) (string, error) {
	return sendCommand(("speed " + speed), d.cmdConn)
}

//StartStreaming command starts video streaming
func (d Drone) StartStreaming() (string, error) {
	return sendCommand("streamon", d.cmdConn)
}

//EndStreaming command ends video streaming
func (d Drone) EndStreaming() (string, error) {
	return sendCommand("streamoff", d.cmdConn)
}

//GetSpeed command gets current speed
func (d Drone) GetSpeed() (string, error) {
	return sendCommand("speed?", d.cmdConn)
}

//GetBattery command gets battery level
func (d Drone) GetBattery() (string, error) {
	return sendCommand("battery?", d.cmdConn)
}

//GetTime command gets current flight time
func (d Drone) GetTime() (string, error) {
	return sendCommand("time?", d.cmdConn)
}

//GetHeight command gets current drone height
func (d Drone) GetHeight() (string, error) {
	return sendCommand("height?", d.cmdConn)
}

//GetTemp command gets the temperature
func (d Drone) GetTemp() (string, error) {
	return sendCommand("temp?", d.cmdConn)
}

//GetAtitude command get Pitch, Roll, Yaw
func (d Drone) GetAtitude() (string, error) {
	return sendCommand("atitude?", d.cmdConn)
}

//GetBaro command gets barometer data
func (d Drone) GetBaro() (string, error) {
	return sendCommand("baro?", d.cmdConn)
}

//GetAcceleration command gets the current acceleration
func (d Drone) GetAcceleration() (string, error) {
	return sendCommand("acceleration?", d.cmdConn)
}

//GetTOF command gets flight distance
func (d Drone) GetTOF() (string, error) {
	return sendCommand("tof?", d.cmdConn)
}

//GetWiFi command gets the signal to noise ratio
func (d Drone) GetWiFi() (string, error) {
	return sendCommand("wifi?", d.cmdConn)
}

func sendCommand(command string, c net.Conn) (string, error) {
	fmt.Println(command)
	_, err := fmt.Fprintf(c, command)
	if err != nil {
		return "none", fmt.Errorf("couldn't send command because %s", err)
	}

	go handleResponse(c)

	return "end", nil
}
