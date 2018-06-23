package vboxmanage

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func runCommand(cmd ...string) (*bufio.Scanner, error) {
	cmdline := strings.Join(cmd, " ")
	logit(0, "Cmd: vboxmanage %s\n", cmdline)
	out, err := exec.Command("vboxmanage", cmd...).CombinedOutput()
	buffer := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(buffer)

	return scanner, err
}
