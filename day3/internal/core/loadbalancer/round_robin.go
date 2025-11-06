package loadbalancer

type RoundRobinLoadBalancer struct {
	*BaseLB
	current int
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		BaseLB:  NewBaseLB(),
		current: -1,
	}
}

// 跳转到下一个可用的后端
func (r *RoundRobinLoadBalancer) GetNextTarget() *Backend {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.backends) == 0 {
		return nil
	}
	cnt := 0
	r.current = (r.current + 1) % len(r.backends)
	for !r.backends[r.current].IsAlive() {
		r.current = (r.current + 1) % len(r.backends)
		cnt++
		if cnt >= len(r.backends) {
			return nil
		}
	}
	return r.backends[r.current]
}
