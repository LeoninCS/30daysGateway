package loadbalancer

import (
	"net/url"
	"sync"
)

type Backend struct {
	ID    string
	URL   *url.URL
	mu    sync.Mutex
	Alive bool
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Alive
}
