package vboxmanage

import (
	"bufio"
	"bytes"
	"os/exec"
)

func runCommand(cmd ...string) (*bufio.Scanner, error) {
	out, err := exec.Command("vboxmanage", cmd...).CombinedOutput()
	buffer := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(buffer)

	return scanner, err
}
