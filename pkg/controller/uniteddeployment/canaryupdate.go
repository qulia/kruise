package uniteddeployment

import "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1"

type CanaryUpdate struct {
	mu  *v1alpha1.ManualUpdate
	dec map[string]int32
}

func NewCanaryUpdate(start *v1alpha1.ManualUpdate, dec map[string]int32) CanaryUpdate {
	cu := CanaryUpdate{mu: start, dec: dec}
	return cu
}
func (cu *CanaryUpdate) Next() {
	if !getNextPartitions(cu.mu.Partitions, cu.dec) {
		cu.mu = nil
	}
}

func (cu CanaryUpdate) Bake() {

}

func (cu CanaryUpdate) Analyze() {

}

func getNextPartitions(partitions map[string]int32, dec map[string]int32) bool {
	set := false
	for key, val := range partitions {
		if val > 0 {
			if dec[key] <= 0 {
				val -= 1
			} else {
				val -= dec[key]
			}
			if val < 0 {
				val = 0
			}

			partitions[key] = val
			set = true
		}
	}

	return set
}
