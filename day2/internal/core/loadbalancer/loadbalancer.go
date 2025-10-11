package loadbalancer

type LoadBalancer interface {
	GetNextTarget() *Backend
	AddBackend(backend *Backend)
	RemoveBackend(backend *Backend)
	GetBackends() []*Backend
}
