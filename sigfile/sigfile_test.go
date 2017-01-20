package sigfile

import (
	"io/ioutil"
	"testing"
)

func TestDecodeSingle(t *testing.T) {
	bytes, err := ioutil.ReadFile("./test-single.sig")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	vals, err := Decode(bytes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(vals) != 1 {
		t.Error("Incorrect number of records")
		t.Fail()
	}

	if vals[0].MemAddr != 22005 {
		t.Error("Bad memory address")
		t.Fail()
	}

	if vals[0].Name != "::PktDeltaLWP:S-1.1:S-3:S-8.10.2.000055F5" {
		t.Error("Bad Signal Name")
		t.Fail()
	}

	if vals[0].SigType[0] != 0 && vals[0].SigType[1] != 1 {
		t.Error("Incorrect type.")
		t.Fail()
	}
}

func TestDecodeFile(t *testing.T) {
	bytes, err := ioutil.ReadFile("./test.sig")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	vals, err := Decode(bytes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(vals) != 19871 {
		t.Error("Incorrect number of records")
		t.Fail()
	}

	lastEntry := vals[len(vals)-1]

	if lastEntry.MemAddr != 5113 {
		t.Error("Bad memory address")
		t.Fail()
	}

	if lastEntry.Name != "Display1_Offline_and" {
		t.Error("Bad Signal Name")
		t.Fail()
	}

	if lastEntry.SigType[0] != 0 || lastEntry.SigType[1] != 1 {
		t.Error("Incorrect type.")
		t.Fail()
	}

	//fmt.Printf("%v\n", len(vals))
	//fmt.Printf("%+v\n", vals[len(vals)-1])
}
