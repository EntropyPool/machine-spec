package machspec

import (
	"errors"
	dmidecode "github.com/dselans/dmidecode"
	"os/user"
)

type MachineSpec struct {
}

func ReadMachineSpec() (*MachineSpec, error) {
	usr, err := user.Current()
	if nil != err {
		return nil, err
	}
	if "root" != usr.Username {
		return nil, errors.New("permission denied: must be run with root")
	}
	dmi := dmidecode.New()
	if err := dmi.Run(); nil != err {
		return nil, err
	}
	return nil, nil
}
