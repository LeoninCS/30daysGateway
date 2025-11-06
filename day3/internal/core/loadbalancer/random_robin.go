package loadbalancer

import (
	"math/rand"
)

type RandomLoadBalancer struct {
	*BaseLB
	current int
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{
		BaseLB:  NewBaseLB(),
		current: -1,
	}
}

// 跳转到下一个可用的后端
func (r *RandomLoadBalancer) GetNextTarget() *Backend {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.backends) == 0 {
		return nil
	}
	cnt := 0
	r.current = rand.Intn(len(r.backends))
	for !r.backends[r.current].IsAlive() {
		r.current = rand.Intn(len(r.backends))
		cnt++
		if cnt >= len(r.backends) {
			return nil
		}
	}
	return r.backends[r.current]
}
