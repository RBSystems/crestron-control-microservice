package helpers

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

var port = ":41795"

func queryState(sigNumber uint32, address string) (string, error) {
	log.Printf("querying state of %v on %v", sigNumber, address)

	tcpAdder, err := net.ResolveTCPAddr("tcp", address+port)

	if err != nil {
		log.Printf("error resolving address. ERROR: %v", err.Error())
		return "", err
	}

	connection, err := net.DialTCP("tcp", nil, tcpAdder)

	if err != nil {
		log.Printf("error connecting to host. ERROR: %v", err.Error())
		return "", err
	}

	defer connection.Close()

	response, err := readUntil(connection, ">")
	if err != nil {
		log.Printf("error reading response. ERROR: %v", err.Error())
		return "", err
	}
	fmt.Printf("%s\n", response)
	fmt.Printf("%v\n", response)

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

	metaResponse, err := readUntil(connection, fmt.Sprintf(`%x=(\S*)\r`, input))
	if err != nil {
		log.Printf("error reading response. ERROR: %v", err.Error())
		return "", err
	}

	//call to readUntil catches any errors
	regEx, _ := regexp.Compile(fmt.Sprintf(`%x=(\S*)\r`, input))
	output := string(regEx.FindSubmatch(metaResponse)[1])

	fmt.Printf("%s\n", output)
	fmt.Printf("%v\n", output)

	return output, nil
}

//sets the state
func setState(sigNumber uint32, sigValue string, address string) error {
	log.Printf("setting state of %v to %v on %v", sigNumber, sigValue, address)
	return nil
}

func readPacket(connection *net.TCPConn) ([]byte, error) {
	response := make([]byte, 1024)
	_, err := connection.Read(response)

	return response, err
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
