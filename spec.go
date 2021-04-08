package machspec

import (
	"errors"
	"fmt"
	dmidecode "github.com/dselans/dmidecode"
	"net"
	"os/user"
	"sort"
	"strconv"
	"strings"
)

type Memory struct {
	Dimm         string `json:"dimm"`
	Manufacturer string `json:"manufacturer"`
	Sn           string `json:"sn"`
	SizeGB       int    `json:"size_gb"`
	Type         string `json:"type"`
	Speed        string `json:"speed"`
}

type MachineSpec struct {
	SerialNumber []map[string]string `json:"sn"` //主板序列号
	Memory       []Memory            `json:"memory"`
	Cpu          []map[string]string `json:"cpu"`
}

func (spec *MachineSpec) MAC() string {
	ifs, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var macs []string
	for _, ifa := range ifs {
		mac := ifa.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, mac)
		}
	}
	sort.SliceStable(macs, func(i, j int) bool {
		return macs[i] < macs[j]
	})
	macsStr := ""
	for i, mac := range macs {
		str := strings.Replace(mac, ":", "", -1)
		if strings.Contains(macsStr, str) {
			continue
		}
		if 0 < i {
			macsStr = fmt.Sprintf("%s-", macsStr)
		}
		macsStr = fmt.Sprintf("%s%s", macsStr, str)
	}
	return macsStr
}

func (spec *MachineSpec) SN() string {
	sn := ""
	keys := make([]string, 0)
	for _, snm := range spec.SerialNumber {
		keys = append(keys, snm["type"])
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	type keyC struct {
		key string
		val string
	}
	keyCs := make([]keyC, len(keys))
	for i, key := range keys {
		keyCs[i].key = key
	}
	for _, snm := range spec.SerialNumber {
		for i, key := range keys {
			if snm["type"] == key {
				keyCs[i].val = snm["serial_number"]
				break
			}
		}
	}
	i := 0
	for _, keyC := range keyCs {
		if 0 < i {
			sn = fmt.Sprintf("%s-", sn)
		}
		sn = fmt.Sprintf("%s%v-%s", sn, keyC.key, keyC.val)
		i += 1
	}
	if "" == sn {
		return spec.MAC()
	}
	return fmt.Sprintf("%s-mac-%s", sn, spec.MAC())
}

func NewMachineSpec() *MachineSpec {
	return &MachineSpec{}
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
func (spec *MachineSpec) PrepareLowLevel() error {
	usr, err := user.Current()
	if nil != err {
		return err
	}
	if "root" != usr.Username {
		return errors.New("permission denied: must be run with root")
	}

	dmi := dmidecode.New()
	if err := dmi.Run(); nil != err {
		return err
	}

	machineSpec := spec

	for _, records := range dmi.Data {
		for _, record := range records {
			//系统序列号
			if record["DMIType"] == "1" || record["DMIType"] == "2" || record["DMIType"] == "3" {
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
				if !strings.Contains(record["Size"], "GB") {
					continue
				}
				tmpMemory := strings.Split(record["Size"], " ")
				size, _ := strconv.ParseInt(tmpMemory[0], 10, 64)

				memory := Memory{}
				memory.SizeGB = int(size)
				memory.Dimm = record["Locator"]
				memory.Manufacturer = record["Manufacturer"]
				memory.Sn = record["Serial Number"]
				memory.Type = record["Type"]
				memory.Speed = record["Speed"]

				machineSpec.Memory = append(machineSpec.Memory, memory)
			}
		}
	}

	return nil
}
