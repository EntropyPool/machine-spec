package machspec

import (
	log "github.com/EntropyPool/entropy-logger"
	"testing"
)

func TestReadMachineSpec(t *testing.T) {
	spec := NewMachineSpec()

	err := spec.PrepareLowLevel()
	if nil != err {
		t.Errorf("CANNOT read machine spec by [%v]", err)
	}
	sn := spec.SN()
	log.Infof(log.Fields{}, "Machine SN: %v", sn)
	macs := spec.MAC()
	log.Infof(log.Fields{}, "Machine MAC: %v", macs)
	log.Infof(log.Fields{}, "Memory: %v", spec.Memory)
}
