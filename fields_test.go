package output

import (
	"testing"
)

func TestFnName(t *testing.T) {
	outputCallerName := FnName()
	if outputCallerName != "TestFnName" {
		t.Fail()
	}
}

func TestCallerName(t *testing.T) {
	instanceCallerName := NewOutputter(nil, nil).CallerName()
	if instanceCallerName != "TestCallerName" {
		t.Fail()
	}
}
