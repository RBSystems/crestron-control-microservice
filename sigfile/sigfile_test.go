package sigfile

import (
	"io/ioutil"
	"testing"
)

func TestDecode(t *testing.T) {
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

	if len(vals) != 10000 {
		t.Fail()
	}
}
