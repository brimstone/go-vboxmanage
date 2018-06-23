package vboxmanage

import "fmt"

func (vm *VM) Poweroff() error {
	scanner, err := runCommand("controlvm", vm.UUID, "poweroff")
	if err != nil {
		return err
	}

	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(s)
		err := vm.Refresh()
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("Unable to start")
}
