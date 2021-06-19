package vboxmanage

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (vm *VM) Refresh() error {
	var snapshot *Snapshot
	longre := regexp.MustCompile(`^"?([^"]*)"?="?([^"]*)"?$`)

	longscanner, err := runCommand("showvminfo", "--machinereadable", vm.UUID)
	if err != nil {
		return fmt.Errorf("Error running `showvminfo`: %s", err)
	}
	for longscanner.Scan() {
		s2 := longscanner.Text()
		results2 := longre.FindStringSubmatch(s2)
		if len(results2) != 3 {
			continue
		}
		if results2[1] == "macaddress1" { // TODO gather all MAC addresses
			if len(vm.MAC) == 0 {
				vm.MAC = append(vm.MAC, "")
			}
			vm.MAC[0] = results2[2][0:2] + ":" +
				results2[2][2:4] + ":" +
				results2[2][4:6] + ":" +
				results2[2][6:8] + ":" +
				results2[2][8:10] + ":" +
				results2[2][10:12]
		} else if results2[1] == "macaddress2" {
			if len(vm.MAC) == 1 {
				vm.MAC = append(vm.MAC, "")
			}
			vm.MAC[1] = results2[2][0:2] + ":" +
				results2[2][2:4] + ":" +
				results2[2][4:6] + ":" +
				results2[2][6:8] + ":" +
				results2[2][8:10] + ":" +
				results2[2][10:12]
		} else if results2[1] == "memory" {
			vm.Memory, err = strconv.Atoi(results2[2])
			if err != nil {
				return fmt.Errorf("Error parsing memory for VM: %s", err)
			}
		} else if results2[1] == "VMState" {
			if results2[2] == "poweroff" {
				vm.Power = "off"
			} else if results2[2] == "poweron" {
				vm.Power = "on"
			} else if results2[2] == "running" {
				vm.Power = "on"
			} else {
				vm.Power = "unknown"
			}
		} else if strings.HasPrefix(results2[1], "SnapshotName") {
			vm.Snapshots = append(vm.Snapshots, Snapshot{
				Name: results2[2],
			})
			snapshot = &vm.Snapshots[len(vm.Snapshots)-1]
		} else if strings.HasPrefix(results2[1], "SnapshotUUID") {
			snapshot.UUID = results2[2]
		} else if results2[1] == "nic1" {
			vm.Nic = results2[2]
		} else if results2[1] == "bridgeadapter1" {
			vm.Bridge = results2[2]
		} else if results2[1] == "groups" {
			vm.Group = results2[2][1:]
		} else if results2[1] == "description" {
			// Parse key=value out of the description
			description := strings.TrimPrefix(s2, "description=")
			unquoted, err := unquote(description)
			// keep looping until we don't have an error about an improperly quoted description
			for err != nil {
				description += "\n"
				longscanner.Scan()
				description += longscanner.Text()
				// unquoted is basically strconv.Unquoted but without a newline check, see below.
				unquoted, err = unquote(description)
			}
			// check each line of the description, now that the quotes are handled
			for _, unquotedLine := range strings.Split(unquoted, "\n") {
				metaMatch := longre.FindStringSubmatch(unquotedLine)
				if len(metaMatch) != 3 {
					continue
				}
				// if the line matches the key=value syntax add it to our Meta attribute
				if vm.Meta == nil {
					vm.Meta = make(map[string]string)
				}
				vm.Meta[metaMatch[1]] = metaMatch[2]
			}
		}
	}
	return nil
}

// stolen from strconv.Unquote
func unquote(s string) (string, error) {
	ErrSyntax := errors.New("Bad syntax")
	n := len(s)
	if n < 2 {
		return "", ErrSyntax
	}
	quote := s[0]
	if quote != s[n-1] {
		return "", ErrSyntax
	}
	s = s[1 : n-1]
	if quote == '`' {
		if contains(s, '`') {
			return "", ErrSyntax
		}
		if contains(s, '\r') {
			// -1 because we know there is at least one \r to remove.
			buf := make([]byte, 0, len(s)-1)
			for i := 0; i < len(s); i++ {
				if s[i] != '\r' {
					buf = append(buf, s[i])
				}
			}
			return string(buf), nil
		}
		return s, nil
	}
	if quote != '"' && quote != '\'' {
		return "", ErrSyntax
	}
	/*
		if contains(s, '\n') {
			return "", ErrSyntax
		}
	*/
	// Is it trivial?  Avoid allocation.
	if !contains(s, '\\') && !contains(s, quote) {
		switch quote {
		case '"':
			return s, nil
		case '\'':
			r, size := utf8.DecodeRuneInString(s)
			if size == len(s) && (r != utf8.RuneError || size != 1) {
				return s, nil
			}
		}
	}
	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.
	for len(s) > 0 {
		c, multibyte, ss, err := strconv.UnquoteChar(s, quote)
		if err != nil {
			return "", err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
		if quote == '\'' && len(s) != 0 {
			// single-quoted must be single character
			return "", ErrSyntax
		}
	}
	return string(buf), nil
}

// contains reports whether the string contains the byte c.
func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}
