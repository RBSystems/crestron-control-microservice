package helpers

import "testing"

var address = "10.6.36.51"

func TestQueryState(t *testing.T) {
	response, err := QueryState(0x02FB, address)
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(response)
}

func TestSetState(t *testing.T) {
	err := SetState(0x0054, "1", address)
	if err != nil {
		t.Log(err.Error())
	}

	err = SetState(0x0054, "0", address)
	if err != nil {
		t.Log(err.Error())
	}
}
