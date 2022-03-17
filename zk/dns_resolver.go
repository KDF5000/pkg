package zk

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/rs/dnscache"
)

// stringShuffle performs a Fisher-Yates shuffle on a slice of strings
func stringShuffle(s []string) {
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// stringDiff will not change the order of elements
func isEqualSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// copy slice to keep origin slice`s order
	aa := make([]string, len(a), len(a))
	bb := make([]string, len(b), len(b))
	copy(aa, a)
	sort.Strings(aa)
	copy(bb, b)
	sort.Strings(bb)

	for i := range aa {
		if aa[i] != bb[i] {
			return false
		}
	}

	return true
}

// DNSHostProvider is a extended version of default HostProvider.
// It will re-query DNS periodically and cache dns result.
type DNSHostProvider struct {
	mu      sync.Mutex // Protects everything, so we can add asynchronous updates later.
	servers []string
	addrs   []string // lookup from servers
	curr    int
	last    int

	resolver   *dnscache.Resolver
	lookupHost func(string) ([]string, error) // Override of net.LookupHost, for testing.
}

// Init is called first, with the servers specified in the connection
// string. It uses DNS to look up addresses for each server, then
// shuffles them all together.
func (hp *DNSHostProvider) Init(servers []string) error {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	hp.resolver = &dnscache.Resolver{}
	hp.servers = servers
	found, err := hp.lookupHostAddr(hp.servers)
	if err != nil || len(found) == 0 {
		return fmt.Errorf("No hosts found for addresses %q, err: %v", servers, err)
	}

	// Randomize the order of the servers to avoid creating hotspots
	stringShuffle(found)

	hp.addrs = found
	hp.curr = -1
	hp.last = -1

	// refresh dns address periodically
	go hp.refresh()

	return nil
}

// Len returns the number of servers available
func (hp *DNSHostProvider) Len() int {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	return len(hp.addrs)
}

// Next returns the next server to connect to. retryStart will be true
// if we've looped through all known servers without Connected() being
// called.
func (hp *DNSHostProvider) Next() (server string, retryStart bool) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.curr = (hp.curr + 1) % len(hp.addrs)
	retryStart = hp.curr == hp.last
	if hp.last == -1 {
		hp.last = 0
	}

	return hp.addrs[hp.curr], retryStart
}

// Connected notifies the HostProvider of a successful connection.
func (hp *DNSHostProvider) Connected() {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.last = hp.curr
}

func (hp *DNSHostProvider) lookupHostAddr(servers []string) ([]string, error) {
	lookupHost := hp.lookupHost
	if lookupHost == nil {
		lookupHost = func(host string) ([]string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return hp.resolver.LookupHost(ctx, host)
		}
	}

	found := []string{}
	for _, server := range servers {
		host, port, err := net.SplitHostPort(server)
		if err != nil {
			return found, err
		}
		addrs, err := lookupHost(host)
		if err != nil {
			return found, err
		}

		for _, addr := range addrs {
			found = append(found, net.JoinHostPort(addr, port))
		}
	}

	return found, nil
}

func (hp *DNSHostProvider) refresh() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()

	for range t.C {
		hp.resolver.Refresh(true)

		// Lookup new addr and update h.addr if there is difference
		addrs, err := hp.lookupHostAddr(hp.servers)
		if err != nil {
			log.Printf("lookupHostAddr error: %v", err)
			continue
		}

		if isEqualSlice(hp.addrs, addrs) {
			continue
		}

		stringShuffle(addrs)

		hp.mu.Lock()
		hp.addrs = addrs
		hp.last = -1
		hp.curr = -1
		hp.mu.Unlock()
	}
}
