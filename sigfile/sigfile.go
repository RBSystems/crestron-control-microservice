package sigfile

import (
	"encoding/binary"
	"log"
)

type Signal struct {
	Name    string
	MemAddr uint32
	SigType []byte
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
