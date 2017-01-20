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

	log.Printf("Sigfile %x", sigfile)

	for pos < len(sigfile) {
		sig := Signal{}

		log.Printf("Pos: %v", pos)

		log.Printf("Size Bytes: %x", sigfile[pos:pos+2])

		//first two bytes are the length
		size := binary.LittleEndian.Uint16(sigfile[pos : pos+2])
		log.Printf("Size: %v", size)

		//the rest composes our signal
		curBytes := sigfile[pos+2 : pos+int(size)]

		//last two bytes are our type
		sig.SigType = curBytes[len(curBytes)-2:]
		log.Printf("Type: %v", sig.SigType)

		//four bytes before that are our memory address
		sig.MemAddr = binary.LittleEndian.Uint32(curBytes[len(curBytes)-6 : len(curBytes)-2])
		log.Printf("Addr Little: %v", sig.MemAddr)

		sig.Name = string(curBytes[:len(curBytes)-6])
		log.Printf("Name: %v", sig.Name)

		toReturn = append(toReturn, sig)
		pos = pos + int(size)
	}
	return toReturn, nil
}
