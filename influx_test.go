package main

import "testing"

func TestSubmitData(t *testing.T) {

	err := SubmitData()
	if err != nil {
		t.Error(err)
	}

}
