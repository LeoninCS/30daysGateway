package loadbalancer

import "sync"

type LoadBalancer interface {
	GetNextTarget() *Backend
	AddBackend(backend *Backend)
	RemoveBackend(backend *Backend)
	GetBackends() []*Backend
}

type BaseLB struct {
	backends []*Backend
	mutex    sync.Mutex
}

func NewBaseLB() *BaseLB {
	return &BaseLB{
		backends: make([]*Backend, 0),
	}
}
func (r *BaseLB) AddBackend(backend *Backend) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.backends = append(r.backends, backend)
}

func (r *BaseLB) RemoveBackend(backend *Backend) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for i, b := range r.backends {
		if b.ID == backend.ID {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			return
		}
	}
}
func (r *BaseLB) GetBackends() []*Backend {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.backends
}
