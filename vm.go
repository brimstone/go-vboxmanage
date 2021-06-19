package vboxmanage

type Snapshot struct {
	Name   string
	UUID   string
	Parent string
}

type VM struct {
	UUID      string
	Bridge    string
	Group     string
	Name      string
	MAC       []string
	Memory    int
	Power     string
	Snapshots []Snapshot
	Nic       string
	Meta      map[string]string
}
