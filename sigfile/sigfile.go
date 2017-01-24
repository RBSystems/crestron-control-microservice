package sigfile

import (
	"archive/zip"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	xmodem "github.com/byuoitav/go-xmodem"
	"github.com/ziutek/telnet"
)

//Signal represents a Creston programming signal
type Signal struct {
	Name    string
	MemAddr uint32
	SigType []byte
}

var sigDirectory = "/tmp/sigfiles/"
var refreshRate = time.Duration(15 * time.Minute)

/*
GetSignalsForAddress returns a current sigfile for the address provided.
*/
func GetSignalsForAddress(address string) (map[string]Signal, error) {
	log.Printf("Getting the signals for address %v", address)
	bytes, err := Read(address)
	if err != nil {
		return nil, err
	}
	toReturn, err := Decode(bytes)

	log.Printf("done")
	return toReturn, err
}

/*
Read checks for a current sig file for the address provided, if none is present,
calls Fetch(). Reads and returns the bytes of the current sig file.

This function may also be used to update the sig file by simply ignoring the bytes
returned.
*/
func Read(address string) ([]byte, error) {
	log.Printf("Checking for a current sig file for %v...", address)

	//Check the sig file present.
	info, err := ioutil.ReadDir(sigDirectory + address)
	if err != nil {
		if !os.IsNotExist(err) {
			//some other error.
			return []byte{}, err
		}
		log.Printf("Directory for %v, did not exist, creating directory...", address)

		//the directory didn't exist, we need to create it, then go get the file.
		err = os.MkdirAll(sigDirectory+address, 0777)
		if err != nil {
			log.Printf("Error creating directory. ERROR: %v", err.Error())
			return []byte{}, err
		}
		log.Printf("Directory created.")

		//Do we need to get the info again to continue in flow below.
		//If not, remove below section.
		info, err = ioutil.ReadDir(sigDirectory + address)

		if err != nil {
			log.Printf("Error: %v", err.Error())

		}
	}

	//There should only be one file
	if len(info) > 1 {
		return []byte{}, errors.New("the sig file directory for " + address + " is malformed. Only one file is permitted in the directory.")
	}

	addr := ""
	//No files, go get it.
	if len(info) == 0 {
		addr, err = Fetch(address)
		if err != nil {
			return []byte{}, err
		}
	} else { // we need to check the mod time.
		if time.Now().Sub(info[0].ModTime()) > refreshRate { // it's been too long, go get it.
			addr, err = Fetch(address)
			if err != nil {
				return []byte{}, err
			}
		} else {
			addr = info[0].Name()
		}
	}
	log.Printf("Reading sig file from time %v...", addr)
	//go read addr, return the bytes
	bytes, err := ioutil.ReadFile(sigDirectory + address + "/" + addr)
	if err != nil {
		return []byte{}, err
	}
	log.Printf("Done.")
	return bytes, err
}

/*
Decode takes a sig file and decodes it
The format is
start
[somechar] :  2 byte length n : record length n-2 : 2 byte length : record length n-2 ...
...
end
where a record is in the format of name (n-8) : memory address (4 byte) : type (2 byte)
*/
func Decode(sigfile []byte) (map[string]Signal, error) {
	log.Printf("Looking for beginning of signals.")
	//start our parse by looking through the array for the close bracket character. 0x5D
	pos := 0
	for pos = 0; pos < len(sigfile); pos++ {
		if sigfile[pos] == 0x5D {
			pos++
			break
		}
	}

	log.Printf("Found start symbol.")

	toReturn := make(map[string]Signal)
	log.Printf("Parsing program signals.")

	for pos < len(sigfile) {
		sig := Signal{}

		//first two bytes are the length
		size := binary.LittleEndian.Uint16(sigfile[pos : pos+2])

		//the rest composes our signal
		curBytes := sigfile[pos+2 : pos+int(size)]

		//last two bytes are our type
		sig.SigType = curBytes[len(curBytes)-2:]

		//four bytes before that are our memory address
		sig.MemAddr = binary.LittleEndian.Uint32(curBytes[len(curBytes)-6 : len(curBytes)-2])

		sig.Name = string(curBytes[:len(curBytes)-6])

		toReturn[sig.Name] = sig

		pos = pos + int(size)
	}

	log.Printf("Found %v signals.", len(toReturn))
	return toReturn, nil
}

//Fetch retrieves a sig file fromthe crestron dump
func Fetch(address string) (string, error) {
	log.Printf("Fetching the sig file for %v", address)

	connection, err := telnet.Dial("tcp", address+":41795")
	if err != nil {
		return "", err
	}
	defer connection.Close()
	connection.SetUnixWriteMode(true)

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	//look for prompt
	output, err := connection.ReadUntil(">")
	if err != nil {
		return "", err
	}
	log.Printf("%s", output)

	connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = connection.Write([]byte("\n"))
	if err != nil {
		return "", err
	}

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	output, err = connection.ReadUntil(">")
	if err != nil {
		return "", err
	}
	log.Printf("%s", output)

	//send command over conncection
	connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = connection.Write([]byte("XGET TEC HD.zig\n"))
	if err != nil {
		return "", err
	}

	log.Print("Read in progress")

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	connection.ReadUntil("DMPS-300-C", "FILE")

	time.Sleep(2 * time.Second)

	log.Print("FILE UPLOAD found")

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	connection.SetWriteDeadline(time.Now().Add(10 * time.Second))

	log.Printf("Starting XMODEM receive")
	response, err := xmodem.Receive(connection.Conn)
	if err != nil {
		return "", err
	}

	log.Printf("XMODEM receive done. ")

	log.Printf("Preparing to unzip files.")
	err = ioutil.WriteFile(sigDirectory+"out-temp.zip", response, 0777)
	if err != nil {
		log.Printf("Could not write temp zip file. ERROR: %v", err.Error())
		return "", err
	}

	r, err := zip.OpenReader(sigDirectory + "out-temp.zip")
	if err != nil {
		log.Printf("error opening zip file for read: ERROR: %v", err.Error())
		return "", err
	}
	defer r.Close()

	log.Printf("Unzipping files.")
	//Validate that the zig only has one file
	if len(r.File) != 1 {
		log.Printf("Zig had more than one file.")
		return "", errors.New("Zig had more than one file.")
	}

	f := r.File[0]
	log.Printf("Writing %s", f.Name)
	rc, err := f.Open()

	if err != nil {
		return "", err
	}

	timestamp := time.Now().Format(time.RFC3339) // The timestamp that acts as the filename

	outfile, er := os.OpenFile(sigDirectory+address+"/"+timestamp+".sig", os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if er != nil {
		return "", err
	}

	_, err = io.Copy(outfile, rc)

	if err != nil {
		return "", err
	}

	rc.Close()
	log.Printf("Extration and inflation done.")

	return timestamp + ".sig", nil

}
