package machspec

import (
	"fmt"
	"strings"
	"errors"
	dmidecode "github.com/dselans/dmidecode"
	"os/user"
	//"encoding/json"
)

type MachineSpec struct {
	SerialNumber []map[string]string `json:"sn"` //主板序列号
	Memory []string `json:"memory"`
	Cpu []map[string]string `json:"cpu"`
}

/**
 * DMIType 编码列表
 * 0 BIOS
 * 1 System
 * 2 Base Bord
 * 3 Chassis
 * 4 Processor
 * 16 Physical Memory Array
 * 17 Memory Device
 */
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

	machineSpec := new(MachineSpec)
	for _,records := range dmi.Data {
		for _, record := range records {
			//系统序列号
			if record["DMIType"] == "1"|| record["DMIType"] == "2" || record["DMIType"] == "3" {
				sn := make(map[string]string)
				sn["type"] = record["DMIType"]
				sn["serial_number"] = record["Serial Number"]
				machineSpec.SerialNumber = append(machineSpec.SerialNumber, sn)
		    }

			//cpu
			if record["DMIType"] == "4" {
				cpu := make(map[string]string)
				cpu["ID"] = record["ID"]
				cpu["Version"] = record["Version"]
				machineSpec.Cpu = append(machineSpec.Cpu, cpu)
			}

			//memory
			if record["DMIType"] == "17" {
				if strings.Contains(record["Size"], "GB") {
					tmpMemory := strings.Split(record["Size"], " ")
					machineSpec.Memory = append(machineSpec.Memory, tmpMemory[0])
				}
			}
		}
	}

	return machineSpec, nil
}
