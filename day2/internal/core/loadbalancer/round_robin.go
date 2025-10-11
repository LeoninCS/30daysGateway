package loadbalancer

import "sync"

type RoundRobinLoadBalancer struct {
	backends []*Backend
	current  int
	mutex    sync.Mutex
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		backends: make([]*Backend, 0),
		current:  -1,
	}
}

// 跳转到下一个可用的后端
func (r *RoundRobinLoadBalancer) GetNextTarget() *Backend {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.backends) == 0 {
		return nil
	}
	r.current = (r.current + 1) % len(r.backends)
	for !r.backends[r.current].IsAlive() {
		r.current = (r.current + 1) % len(r.backends)
	}
	return r.backends[r.current]
}

func (r *RoundRobinLoadBalancer) AddBackend(backend *Backend) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.backends = append(r.backends, backend)
}

func (r *RoundRobinLoadBalancer) RemoveBackend(backend *Backend) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for i, b := range r.backends {
		if b.ID == backend.ID {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			return
		}
	}
}

func (r *RoundRobinLoadBalancer) GetBackends() []*Backend {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.backends
}
