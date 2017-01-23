package helpers

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var port = ":41795"

func QueryState(sigNumber uint32, address string) (string, error) {
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

	response, err = readUntil(connection, "0000")

	if err != nil {
		log.Printf("error reading response. ERROR: %v", err.Error())
		return "", err
	}

	fmt.Printf("%s\n", response)

	return string(response), nil
}

//sets the state
func SetState(sigNumber uint32, sigValue string, address string) error {
	log.Printf("setting state of %v to %v on %v", sigNumber, sigValue, address)
	return nil
}

func readPacket(connection *net.TCPConn) ([]byte, error) {
	response := make([]byte, 1024)
	_, err := connection.Read(response)

	return response, err
}

func readUntil(connection *net.TCPConn, delim string) ([]byte, error) {
	size := len(delim)
	c := make([]byte, size)
	toReturn := []byte{}
	for !strings.Contains(string(toReturn), delim) {

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
