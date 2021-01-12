package machspec

import (
	"errors"
	dmidecode "github.com/dselans/dmidecode"
)

type MachineSpec struct {
}

func ReadMachineSpec() (*MachineSpec, error) {
	dmi := dmidecode.New()
	if err := dmi.Run(); nil != err {
		return nil, errors.New("Fail to run dmidecode")
	}
	return nil, nil
}
