package main

import (
	"strings"
	"sync"

	"github.com/canonical/lxd/lxd/events"
)

// A Daemon can respond to requests from a shared client.
type Daemon struct {
	// Event servers
	events *events.Server

	// ContextID and port of the LXD VM socket server.
	serverCID         uint32
	serverPort        uint32
	serverCertificate string

	// The channel which is used to indicate that the lxd-agent was able to connect to LXD.
	chConnected chan struct{}

	devlxdRunning bool
	devlxdMu      sync.Mutex
	devlxdEnabled bool

	// Network interfaces to exclude from stats/state queries
	excludedInterfaces map[string]bool
}

// newDaemon returns a new Daemon object with the given configuration.
func newDaemon(debug, verbose bool, excludeInterfaces string) *Daemon {
	lxdEvents := events.NewServer(debug, verbose, nil)

	// Parse comma-separated list of interfaces to exclude
	excludedInterfaces := make(map[string]bool)
	if excludeInterfaces != "" {
		for _, iface := range strings.Split(excludeInterfaces, ",") {
			name := strings.TrimSpace(iface)
			if name != "" {
				excludedInterfaces[name] = true
			}
		}
	}

	return &Daemon{
		events:             lxdEvents,
		chConnected:        make(chan struct{}),
		excludedInterfaces: excludedInterfaces,
	}
}
