package vboxmanage

import (
	"errors"
	"fmt"
	"strings"
)

func (vm *VM) RestoreSnapshot(name string) error {
	snapshotExists := false
	for _, s := range vm.Snapshots {
		if s.Name == name {
			snapshotExists = true
		}
	}
	if !snapshotExists {
		return errors.New("Snapshot doesn't exist")
	}
	runCommand("controlvm", vm.UUID, "poweroff")
	scanner, err := runCommand("snapshot", vm.UUID, "restore", name)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasSuffix(s, "100%") {
			err := vm.Refresh()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("Unable to restore snapshot")
}
