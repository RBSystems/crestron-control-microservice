package helpers

import "testing"

func TestQueryState(t *testing.T) {
	response, err := QueryState(0x02FB, "10.6.36.220")
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(response)
}
