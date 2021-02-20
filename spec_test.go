package machspec

import (
	log "github.com/EntropyPool/entropy-logger"
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
	if nil != spec {
		sn := spec.SN()
		log.Infof(log.Fields{}, "Machine SN: %v", sn)
	}
}
