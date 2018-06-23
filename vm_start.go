package vboxmanage

import "fmt"

func (vm *VM) Start() error {
	scanner, err := runCommand("startvm", vm.UUID, "--type", "headless")
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
