package sigfile

import (
	"archive/zip"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
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

//Log will correspond to the log of when we last checked for the latest programming.
type Log struct {
	Log map[string]time.Time
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

func getLog() (Log, error) {
	var l Log
	l.Log = make(map[string]time.Time)

	logBytes, err := ioutil.ReadFile(sigDirectory + "log.json")
	if err != nil {
		if !os.IsNotExist(err) {
			//some other error.
			return l, err
		}
		log.Printf("Log file does not exist, creating.")

		err = os.MkdirAll(sigDirectory, 0777)
		if err != nil {
			return l, err
		}
		_, err = os.Create(sigDirectory + "log.json")
		if err != nil {
			log.Printf("Problem creating the file %s", err.Error())
			return Log{}, err
		}
		return l, nil
	}

	if len(logBytes) == 0 {
		return l, nil
	}

	err = json.Unmarshal(logBytes, &l)
	if err != nil {
		return l, err
	}

	return l, nil
}

//This might be an issue if we're writing to the log fairly often.
func writeLog(l Log) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(sigDirectory+"log.json", b, 0777)
	if err != nil {
		return err
	}
	return nil
}

/*
Read checks for a current sig file for the address provided, if none is present,
calls Fetch(). Reads and returns the bytes of the current sig file.

This function may also be used to update the sig file by simply ignoring the bytes
returned.
*/
func Read(address string) ([]byte, error) {
	log.Printf("Checking for a current sig file for %v...", address)

	//Check our log to see if we've checked it within the timeframe.
	l, er := getLog()

	if er != nil {
		return []byte{}, er
	}
	addr := ""

	//check the log file for the address
	if date, ok := l.Log[address]; ok {
		if time.Now().Sub(date) > refreshRate {
			//we need to go check for it.

			info, err := getSingleFileFromDirectory(sigDirectory + address)
			if err != nil {
				return []byte{}, err
			}

			//No files, go get it.
			if info == nil {
				addr, err = Fetch(address)
				if err != nil {
					return []byte{}, err
				}
			} else { // we need to check the mod time.
				t, err := GetCompileTime(address)
				if err != nil {
					return []byte{}, err
				}

				fileTime, err := time.Parse(time.RFC3339, strings.Split(info.Name(), ".sig")[0])
				if err != nil {
					return []byte{}, err
				}

				if t.Sub(fileTime) != 0 { // compile times don't match

					//first remove the old file
					err = os.Remove(sigDirectory + address + "/" + info.Name())
					if err != nil {
						log.Printf("Error deleting old sig file. ERROR: %v", err.Error())
						return []byte{}, err
					}

					//fetch the new one
					addr, err = Fetch(address)
					fmt.Printf("Addr0: %v\n", addr)
					if err != nil {
						return []byte{}, err
					}

					//Mark the last time we checked.
					l.Log[address] = time.Now()
				} else {

					addr = info.Name()
					fmt.Printf("Addr1: %v\n", addr)
					//Mark the last time we checked.
					l.Log[address] = time.Now()
				}
			}
		} else { //the file exists, and is within the refresh rate. Get the name of it so we can retrieve it below.

			singleInfo, err := getSingleFileFromDirectory(sigDirectory + address)
			if err != nil {
				return []byte{}, err
			}
			addr = singleInfo.Name()

			l.Log[address] = time.Now()
		}
	} else { //There isn't an entry for this address in the logging map, therefore we need to go get it and add a directory.
		var err error
		//not there, we need to go get it.
		addr, err = Fetch(address)
		if err != nil {
			return []byte{}, err
		}

		fmt.Printf("Addr2: %v\n", addr)
		l.Log[address] = time.Now()
		//add the entry to the log
	}

	fmt.Printf("Addr3: %v\n", addr)

	log.Printf("Reading sig file from time %v...", addr)
	//go read addr, return the bytes
	bytes, err := ioutil.ReadFile(sigDirectory + address + "/" + addr)
	if err != nil {
		return []byte{}, err
	}
	log.Printf("Done.")

	//save out the log file
	writeLog(l)
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

		toReturn[strings.ToLower(sig.Name)] = sig

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
	err = sendNewlineWaitForPrompt(connection)
	if err != nil {
		return "", err
	}

	compileTime, err := getCompileTimeFromConnection(connection)
	if err != nil {
		return "", err
	}

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	output, err := connection.ReadUntil(">")
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
	output, err = connection.ReadUntil("DMPS-300-C", "FILE")
	if err != nil {
		return "", err
	}
	log.Printf("%v", output)

	time.Sleep(2 * time.Second)

	log.Print("FILE UPLOAD found")

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	connection.SetWriteDeadline(time.Now().Add(10 * time.Second))

	log.Printf("Starting XMODEM receive")
	response, err := xmodem.Receive(connection.Conn)
	if err != nil {
		return "", err
	}

	fmt.Printf("response length: %v\n", len(response))

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
	defer rc.Close()

	if err != nil {
		return "", err
	}

	timestamp := compileTime.Format(time.RFC3339) // The timestamp that acts as the filename

	//Open the file.
	outfile, err := os.OpenFile(sigDirectory+address+"/"+timestamp+".sig", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil { //if we get a ISNotExist error that means the directory doesn't exist, we need to create it.
		if !os.IsNotExist(err) {
			//some other error.
			return "", err
		}
		log.Printf("Directory %v does not exists, creating.", sigDirectory+address)

		err = os.MkdirAll(sigDirectory+address, 0777) //Try creating the directory.
		if err != nil {
			return "", err
		}
		outfile, err = os.OpenFile(sigDirectory+address+"/"+timestamp+".sig", os.O_CREATE|os.O_WRONLY, 0777) //try creating/opening the file again.
		if err != nil {
			return "", err
		}
	}

	_, err = io.Copy(outfile, rc)

	if err != nil {
		return "", err
	}

	log.Printf("timestamp: %v", timestamp)

	log.Printf("Extration and inflation done.")

	return timestamp + ".sig", nil

}

func getCompileTimeFromConnection(connection *telnet.Conn) (time.Time, error) {
	log.Printf("Getting the compile time...")

	err := sendNewlineWaitForPrompt(connection)

	connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = connection.Write([]byte("progcom\n"))
	if err != nil {
		return time.Time{}, err
	}

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	output, err := connection.ReadUntil("Compiler Rev")
	if err != nil {
		return time.Time{}, err
	}
	//log.Printf("%s", output)

	//find the line that says Compiled On: .....\n

	regexString := `Compiled On: (.*)[\n\r]+`
	regEx, _ := regexp.Compile(regexString)
	dateString := strings.TrimSpace(string(regEx.FindSubmatch(output)[1]))

	date, err := time.Parse("1/2/2006 3:04 PM", dateString)
	if err != nil {
		return time.Time{}, err
	}

	log.Printf("Compiled on %v.", date.Format(time.RFC3339))
	return date, nil
}

func GetCompileTime(address string) (time.Time, error) {
	log.Printf("Checking program compile time for %s", address)
	connection, err := telnet.Dial("tcp", address+":41795")
	if err != nil {
		return time.Time{}, err
	}
	defer connection.Close()
	connection.SetUnixWriteMode(true)

	return getCompileTimeFromConnection(connection)
}
