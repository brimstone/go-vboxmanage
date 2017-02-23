package vboxmanage

import (
	"fmt"
	"regexp"
)

func ListVMs() ([]VM, error) {
	var vms []VM
	listre := regexp.MustCompile(`^"(.*)" {(.*)}$`)
	longre := regexp.MustCompile(`^"?(.*)"?="?([^"]*)"$`)

	scanner, err := runCommand("list", "vms")
	if err != nil {
		return []VM{}, fmt.Errorf("Error running `list vms`: %s", err)
	}

	for scanner.Scan() {
		s := scanner.Text()
		results := listre.FindStringSubmatch(s)

		if len(results) != 3 {
			continue
		}

		vm := VM{
			Name: results[1],
			UUID: results[2],
		}

		longscanner, err := runCommand("showvminfo", "--machinereadable", vm.UUID)
		if err != nil {
			return []VM{}, fmt.Errorf("Error running `showvminfo`: %s", err)
		}
		for longscanner.Scan() {
			s2 := longscanner.Text()
			results2 := longre.FindStringSubmatch(s2)
			if len(results2) != 3 {
				continue
			}
			if results2[1] == "macaddress1" {
				vm.MAC = results2[2]
			}
		}

		vms = append(vms, vm)
	}

	return vms, nil
}
