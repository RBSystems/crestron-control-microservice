package sigfile

import (
	"archive/zip"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ziutek/telnet"
)

type Signal struct {
	Name    string
	MemAddr uint32
	SigType []byte
}

func Read(address string) ([]bytes, error) {

}

/*
  Decode takes a sig file and decodes it
  The format is
  start
  [somechar] :  2 byte length n : record length n-2 : 2 byte length : record length n-2 ...
  ...

  where a record is in the format of name (n-8) : memory address (4 byte) : type (2 byte)
*/
func Decode(sigfile []byte) ([]Signal, error) {
	log.Printf("Looking for beginning of signals.")
	//start our parse by looking through the array looking for the close bracket character. 0x5D
	pos := 0
	for pos = 0; pos < len(sigfile); pos++ {
		if sigfile[pos] == 0x5D {
			pos++
			break
		}
	}

	log.Printf("Found start symbol.")

	toReturn := []Signal{}
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

		toReturn = append(toReturn, sig)
		pos = pos + int(size)
	}

	log.Printf("Found %v signals.", len(toReturn))
	return toReturn, nil
}

func Fetch(address string) (string, error) {
	connection, err := telnet.Dial("tcp", address+":41795")
	if err != nil {
		return "", err
	}
	defer connection.Close()
	connection.SetUnixWriteMode(true)

	//look for prompt
	output, err := connection.ReadUntil(">")
	if err != nil {
		return "", err
	}
	log.Printf("%s", output)

	_, err = connection.Write([]byte("\n"))
	if err != nil {
		return "", err
	}

	output, err = connection.ReadUntil(">")
	if err != nil {
		return "", err
	}
	log.Printf("%s", output)

	//send command over conncection
	_, err = connection.Write([]byte("XGET TEC HD.zig\n"))
	if err != nil {
		return "", err
	}

	log.Print("Read in progress\n")

	connection.ReadUntil("DMPS-300-C", "FILE")

	time.Sleep(3 * time.Second)

	log.Print("FILE UPLOAD found")

	response, err := xmodem.Receive(connection.Conn)

	if err != nil {
		return "", err
	}

	ioutil.WriteFile("./out-temp.zip", response, 0777)
	r, err := zip.OpenReader("./out-temp.zip")
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		log.Printf("Writing %s", f.Name)
		rc, err := f.Open()

		if err != nil {
			return "", err
		}

		outfile, er := os.OpenFile("/tmp/sigfiles/"+address+"/"+time.Now().Format(time.RFC3339)+".sig", os.O_CREATE|os.O_WRONLY, os.ModeAppend)
		if er != nil {
			return "", err
		}

		_, err = io.Copy(outfile, rc)

		if err != nil {
			return "", err
		}

		rc.Close()
	}
}
