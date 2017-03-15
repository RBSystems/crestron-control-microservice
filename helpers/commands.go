package helpers

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

var port = ":41795"

func QueryState(sigNumber uint32, address string) (string, error) {
	log.Printf("querying state of %v on %v", sigNumber, address)

	connection, err := startConnection(address, port)
	if err != nil {
		return "", err
	}
	defer connection.Close()

	err = writeBytes(connection, []byte(fmt.Sprintf("DBGSIGNAL %v ON\r\n", sigNumber)))
	if err != nil {
		log.Printf("error writing to connection. ERROR: %v", err.Error())
		return "", err
	}

	err = writeBytes(connection, []byte(fmt.Sprintf("MDBGSIGNAL -S:SYNC %v\r\n", sigNumber)))
	if err != nil {
		log.Printf("error writing to connection. ERROR: %v", err.Error())
		return "", err
	}

	input := make([]byte, 4)
	binary.BigEndian.PutUint32(input, sigNumber)

	regExString := fmt.Sprintf(`%X=(\S*)\r`, input)

	metaResponse, err := readUntil(connection, regExString)
	if err != nil {
		log.Printf("error reading response. ERROR: %v", err.Error())
		return "", err
	}

	//call to readUntil catches any errors
	regEx, _ := regexp.Compile(regExString)
	output := string(regEx.FindSubmatch(metaResponse)[1])

	//if it contains a bracket, it's a hex representation of a byte array
	if strings.Contains(output, "[") {
		//remove first and last character
		output = output[1 : len(output)-1]

		//split on brackets
		elements := strings.Split(output, "][")

		//join strings
		output = strings.Join(elements, "")

		outputBytes, err := hex.DecodeString(output)
		if err != nil {
			return "", err
		}

		output = string(outputBytes)

	}

	return output, nil
}

//sets the state
func SetState(sigNumber uint32, sigValue string, address string) error {
	log.Printf("setting state of %v to %v on %v", sigNumber, sigValue, address)

	connection, err := startConnection(address, port)
	if err != nil {
		return err
	}
	defer connection.Close()

	payload := []byte(fmt.Sprintf("SETSIGNAL %v %v\r\n", sigNumber, sigValue))

	err = writeBytes(connection, payload)
	if err != nil {
		return err
	}

	response, err := QueryState(sigNumber, address)
	if err != nil {
		log.Printf("There was an error setting the state of %s", sigNumber)
	}

	if !strings.EqualFold(response, sigValue) {
		log.Printf("There was a problem setting the value. Value came back as %s, attempting to set to %s", response, sigValue)
	}

	return nil
}

func readPacket(connection *net.TCPConn) ([]byte, error) {
	response := make([]byte, 1024)
	_, err := connection.Read(response)

	return response, err
}

//opens connection, performs handshake, waits for first prompt
func startConnection(address string, port string) (*net.TCPConn, error) {
	tcpAdder, err := net.ResolveTCPAddr("tcp", address+port)
	if err != nil {
		log.Printf("error resolving address. ERROR: %v", err.Error())
		return nil, err
	}

	connection, err := net.DialTCP("tcp", nil, tcpAdder)
	if err != nil {
		log.Printf("error connecting to host. ERROR: %v", err.Error())
		return nil, err
	}

	_, err = readUntil(connection, ">")
	if err != nil {
		log.Printf("error reading response. ERROR: %v", err.Error())
		return nil, err
	}

	return connection, nil
}

func readUntil(connection *net.TCPConn, expression string) ([]byte, error) {
	size := len(expression)
	c := make([]byte, size)
	toReturn := []byte{}

	regEx, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	for !regEx.Match(toReturn) {

		_, err := connection.Read(c)
		if err != nil {
			return []byte{}, err
			log.Printf("Error: %s", err.Error())
		}
		toReturn = append(toReturn, c...)
	}

	return toReturn, nil
}

func writeBytes(connection *net.TCPConn, payload []byte) error {
	err := connection.SetWriteDeadline(time.Now().Add((10 * time.Second)))
	if err != nil {
		return err
	}
	_, err = connection.Write(payload)

	return err
}
