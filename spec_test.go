package machspec

import (
	"testing"
)

func TestReadMachineSpec(t *testing.T) {
	spec, err := ReadMachineSpec()
	if nil != err {
		t.Errorf("CANNOT read machine spec by [%v]", err)
	}
	if nil == spec {
		t.Errorf("CANNOT read machine spec by [return nil]")
	}
}
