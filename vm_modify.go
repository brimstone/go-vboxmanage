package vboxmanage

import (
	"errors"
	"strings"
)

func (vm *VM) Modify(key string, value string) error {
	cmd := []string{
		"modifyvm",
		vm.UUID,
		"--" + key,
		value,
	}
	_, err := runCommand(cmd...)
	if err != nil {
		return errors.New("Error running cmd: " + strings.Join(cmd, " ") + " " + err.Error())
	}
	err = vm.Refresh()
	if err != nil {
		return err
	}
	return nil
}
