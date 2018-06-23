package vboxmanage

import (
	"fmt"
	"regexp"
)

func ListVMs() ([]VM, error) {
	var vms []VM
	listre := regexp.MustCompile(`^"(.*)" {(.*)}$`)

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

		err := vm.Refresh()
		if err != nil {
			return nil, err
		}

		vms = append(vms, vm)
	}

	return vms, nil
}
