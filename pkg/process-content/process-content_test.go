package process_content

import (
	"testing"
)

func TestHasContentChanged(t *testing.T){
	a := "123"
	b := "456"
	c := "123"

	abResult := hasContentChanged(a,b)
	if abResult == true {
		t.Errorf("hasContentChanged, has %t, expected: %t", abResult, false)
	}

	acResult := hasContentChanged(a,c)
	if acResult == false {
		t.Errorf("hasContentChanged, has %t, expected: %t", abResult, false)
	}

}