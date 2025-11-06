package loadbalancer

type LeastConnLoadBalancer struct {
	*BaseLB
	currenttime map[*Backend]int
}

func NewLeastConnLoadBalancer() *LeastConnLoadBalancer {
	return &LeastConnLoadBalancer{
		BaseLB:      NewBaseLB(),
		currenttime: make(map[*Backend]int),
	}
}

func (w *LeastConnLoadBalancer) GetNextTarget() *Backend {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.backends) == 0 {
		return nil
	}
	idx, minConn := -1, int(^uint(0)>>1)
	for i, backend := range w.backends {
		if !backend.IsAlive() {
			continue
		}
		if w.currenttime[backend] < minConn {
			minConn = w.currenttime[backend]
			idx = i
		}
	}
	if idx == -1 {
		return nil
	}
	w.currenttime[w.backends[idx]]++
	return w.backends[idx]
}
