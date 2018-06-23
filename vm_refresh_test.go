package vboxmanage_test

import (
	"log"
	"testing"

	"github.com/brimstone/go-vboxmanage"
)

func TestRefresh(t *testing.T) {
	vms, err := vboxmanage.ListVMs()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	for _, vm := range vms {
		log.Printf("%#v\n", vm)
	}
}
