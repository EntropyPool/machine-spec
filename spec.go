package machspec

import (
	"fmt"
	"errors"
	dmidecode "github.com/dselans/dmidecode"
	"os/user"
	//"encoding/json"
)

type MachineSpec struct {
	SerialNumber []string `json:"sn"` //主板序列号
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

	fmt.Println(dmi.Data);
	fmt.Println("--------------")

	machineSpec := new(MachineSpec)
	for _,records := range dmi.Data {
		for _, record := range records {
			for key, val := range record {
				if key == "Serial Number" {
					machineSpec.SerialNumber = append(machineSpec.SerialNumber, val)
				}
			}
		}
	}

	return machineSpec, nil
}
