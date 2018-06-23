package vboxmanage

import (
	"fmt"
	"strings"
)

func (vm *VM) MakeSnapshot(name string, description string) error {
	scanner, err := runCommand("snapshot", vm.UUID, "take", name)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "Snapshot taken") {
			err := vm.Refresh()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("Unable to take snapshot")
}
