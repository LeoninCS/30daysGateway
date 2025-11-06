package loadbalancer

import (
	"math"
)

type WeightedLoadBalancer struct {
	*BaseLB
	currentweight map[*Backend]int
}

func NewWeightedLoadBalancer() *WeightedLoadBalancer {
	return &WeightedLoadBalancer{
		BaseLB:        NewBaseLB(),
		currentweight: make(map[*Backend]int),
	}
}

func (w *WeightedLoadBalancer) GetNextTarget() *Backend {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.backends) == 0 {
		return nil
	}
	idx, maxWeight, sumWeight := -1, math.MinInt32, 0
	for i, backend := range w.backends {
		if !backend.IsAlive() {
			continue
		}
		w.currentweight[backend] += backend.Weight
		sumWeight += backend.Weight
		if w.currentweight[backend] > maxWeight {
			maxWeight = w.currentweight[backend]
			idx = i
		}
	}
	if idx == -1 {
		return nil
	}
	w.currentweight[w.backends[idx]] -= sumWeight
	return w.backends[idx]
}
