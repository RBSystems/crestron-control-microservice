package sigfile

import "encoding/binary"

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

	//start our parse by looking through the array looking for the close bracket character. 0x5D
	pos := 0
	for pos = 0; pos < len(sigfile); pos++ {
		if sigfile[pos] == 0x5D {
			break
		}
	}

	toReturn := []Signal{}

	for pos < len(sigfile) {
		sig := Signal{}

		//first two bytes are the length
		size := binary.BigEndian.Uint16(sigfile[pos : pos+2])

		//the rest composes our signal
		curBytes := sigfile[pos+2 : pos+int(size)]

		//last two bytes are our type
		sig.SigType = curBytes[:len(curBytes)-2]

		//eight bytes before that are our memory address
		sig.MemAddr = binary.LittleEndian.Uint32(curBytes[len(curBytes)-10 : len(curBytes)-2])

		sig.Name = string(curBytes[:len(curBytes)-10])

		toReturn = append(toReturn, sig)
	}
	return toReturn, nil
}
