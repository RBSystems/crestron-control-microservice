package sigfile

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ziutek/telnet"
)

func readOrCreateDirectory(directory string) ([]os.FileInfo, error) {
	//Check the sig file present, grab the compile time.
	info, err := ioutil.ReadDir(directory)
	if err != nil {
		if !os.IsNotExist(err) {
			//some other error.
			return []os.FileInfo{}, err
		}
		log.Printf("Directory %v, did not exist, creating directory...", directory)

		//the directory didn't exist, we need to create it, then go get the file.
		err = os.MkdirAll(directory, 0777)
		if err != nil {
			log.Printf("Error creating directory. ERROR: %v", err.Error())
			return []os.FileInfo{}, err
		}
		log.Printf("Directory created.")

		//Do we need to get the info again to continue in flow below.
		//If not, remove below section.
		info, err = ioutil.ReadDir(directory)

		if err != nil {
			log.Printf("Error: %v", err.Error())
			return []os.FileInfo{}, err
		}
	}
	return info, nil
}

func getSingleFileFromDirectory(directory string) (os.FileInfo, error) {
	info, err := readOrCreateDirectory(directory)
	if err != nil {
		return nil, err
	}

	if len(info) > 1 {
		return nil, errors.New("the sig file directory for " + directory + " is malformed. Only one file is permitted in the directory.")
	}
	if len(info) < 1 {
		return nil, nil
	}

	return info[0], nil
}

func sendNewlineWaitForPrompt(connection *telnet.Conn) error {
	connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := connection.Write([]byte("\n"))
	if err != nil {
		return err
	}

	connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	output, err := connection.ReadUntil(">")
	if err != nil {
		return err
	}
	log.Printf("%s", output)
	return nil
}
