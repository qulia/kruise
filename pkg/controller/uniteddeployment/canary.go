package uniteddeployment

import (
	"fmt"
	"time"

	appsv1alpha1 "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func rollPartitions(r *ReconcileUnitedDeployment, instance *appsv1alpha1.UnitedDeployment, nameToSubset *map[string]*Subset,
	nextReplicas *map[string]int32, nextPartitions *map[string]int32, currentRevision *appsv1.ControllerRevision,
	updatedRevision *appsv1.ControllerRevision, subsetType subSetType,
	oldStatus *appsv1alpha1.UnitedDeploymentStatus, collisionCount int32, control ControlInterface, expectedRevision string) (reconcile.Result, error) {
	var newStatus *appsv1alpha1.UnitedDeploymentStatus
	for {
		var err error
		newStatus, err = r.manageSubsets(instance, nameToSubset, nextReplicas, nextPartitions, currentRevision, updatedRevision, subsetType)
		if err != nil {
			klog.Errorf("Fail to update UnitedDeployment %s/%s: %s", instance.Namespace, instance.Name, err)
			r.recorder.Event(instance.DeepCopy(), corev1.EventTypeWarning, fmt.Sprintf("Failed%s", eventTypeSubsetsUpdate), err.Error())
		}
		if !checkAnalysis() {
			return reconcile.Result{}, fmt.Errorf("failed run")
		}

		if !getNextPartitions(*nextPartitions, instance) {
			break
		}
		time.Sleep(time.Duration(instance.Spec.UpdateStrategy.CanaryUpdate.BakeTimeSeconds) * time.Second)
	}

	return r.updateStatus(instance, newStatus, oldStatus, nameToSubset, nextReplicas, nextPartitions,
		currentRevision, updatedRevision, collisionCount, control)
}

func checkAnalysis() bool {
	return true
}

func getNextPartitions(partitions map[string]int32, instance *appsv1alpha1.UnitedDeployment) bool {
	dec := instance.Spec.UpdateStrategy.CanaryUpdate.RollCount
	set := false
	for key, val := range partitions {
		if val > 0 {
			val -= dec
			if val < 0 {
				val = 0
			}

			partitions[key] = val
			set = true
		}
	}

	return set
}
