package vboxmanage

import "log"

var Loglevel = 3

func logit(level int, format string, contents ...string) {
	if level >= Loglevel {
		log.Printf(format, contents)
	}
}
