package domain

import (
	"net"
	"time"
)

// HostsGenerator generates hosts entries based on the provided blockers.
type HostsGenerator struct {
	blockers []Blocker
}

func NewHostsGenerator(blockers []Blocker) *HostsGenerator {
	return &HostsGenerator{blockers: blockers}
}

type HostsEntry struct {
	IP     net.IP
	Domain string
}

func (e HostsEntry) String() string {
	return e.IP.String() + " " + e.Domain
}

func (g *HostsGenerator) Gen(t time.Time) []HostsEntry {
	var entries []HostsEntry
	for _, blocker := range g.blockers {
		if blocker.IsBlocked(t) {
			entries = append(entries, HostsEntry{
				IP:     blocker.ForwardTo,
				Domain: blocker.Domain,
			})
		}
	}
	return entries
}
