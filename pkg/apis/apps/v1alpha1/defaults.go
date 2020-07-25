/*
Copyright 2019 The Kruise Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/openkruise/kruise/pkg/webhook/default_server/utils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "k8s.io/kubernetes/pkg/apis/core/v1"
	utilpointer "k8s.io/utils/pointer"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_SidecarSet set default values for SidecarSet.
func SetDefaults_SidecarSet(obj *SidecarSet) {
	setSidecarSetUpdateStratety(&obj.Spec.Strategy)

	for i := range obj.Spec.Containers {
		setSidecarDefaultContainer(&obj.Spec.Containers[i])
	}
}

func setSidecarSetUpdateStratety(strategy *SidecarSetUpdateStrategy) {
	if strategy.RollingUpdate == nil {
		rollingUpdate := RollingUpdateSidecarSet{}
		strategy.RollingUpdate = &rollingUpdate
	}
	if strategy.RollingUpdate.MaxUnavailable == nil {
		maxUnavailable := intstr.FromInt(1)
		strategy.RollingUpdate.MaxUnavailable = &maxUnavailable
	}
}

func setSidecarDefaultContainer(sidecarContainer *SidecarContainer) {
	container := &sidecarContainer.Container
	v1.SetDefaults_Container(container)
	for i := range container.Ports {
		p := &container.Ports[i]
		v1.SetDefaults_ContainerPort(p)
	}
	for i := range container.Env {
		e := &container.Env[i]
		if e.ValueFrom != nil {
			if e.ValueFrom.FieldRef != nil {
				v1.SetDefaults_ObjectFieldSelector(e.ValueFrom.FieldRef)
			}
		}
	}
	v1.SetDefaults_ResourceList(&container.Resources.Limits)
	v1.SetDefaults_ResourceList(&container.Resources.Requests)
	if container.LivenessProbe != nil {
		v1.SetDefaults_Probe(container.LivenessProbe)
		if container.LivenessProbe.Handler.HTTPGet != nil {
			v1.SetDefaults_HTTPGetAction(container.LivenessProbe.Handler.HTTPGet)
		}
	}
	if container.ReadinessProbe != nil {
		v1.SetDefaults_Probe(container.ReadinessProbe)
		if container.ReadinessProbe.Handler.HTTPGet != nil {
			v1.SetDefaults_HTTPGetAction(container.ReadinessProbe.Handler.HTTPGet)
		}
	}
	if container.Lifecycle != nil {
		if container.Lifecycle.PostStart != nil {
			if container.Lifecycle.PostStart.HTTPGet != nil {
				v1.SetDefaults_HTTPGetAction(container.Lifecycle.PostStart.HTTPGet)
			}
		}
		if container.Lifecycle.PreStop != nil {
			if container.Lifecycle.PreStop.HTTPGet != nil {
				v1.SetDefaults_HTTPGetAction(container.Lifecycle.PreStop.HTTPGet)
			}
		}
	}
}

// SetDefaults_BroadcastJob set default values for BroadcastJob.
func SetDefaults_BroadcastJob(obj *BroadcastJob) {
	utils.SetDefaultPodTemplate(&obj.Spec.Template.Spec)
	if obj.Spec.CompletionPolicy.Type == "" {
		obj.Spec.CompletionPolicy.Type = Always
	}

	if obj.Spec.Parallelism == nil {
		parallelism := int32(1<<31 - 1)
		parallelismIntStr := intstr.FromInt(int(parallelism))
		obj.Spec.Parallelism = &parallelismIntStr
	}

	if obj.Spec.FailurePolicy.Type == "" {
		obj.Spec.FailurePolicy.Type = FailurePolicyTypeFailFast
	}
}

// SetDefaults_StatefulSet set default values for StatefulSet.
func SetDefaults_StatefulSet(obj *StatefulSet) {
	if len(obj.Spec.PodManagementPolicy) == 0 {
		obj.Spec.PodManagementPolicy = appsv1.OrderedReadyPodManagement
	}

	if obj.Spec.UpdateStrategy.Type == "" {
		obj.Spec.UpdateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType

		// UpdateStrategy.RollingUpdate will take default values below.
		obj.Spec.UpdateStrategy.RollingUpdate = &RollingUpdateStatefulSetStrategy{}
	}

	if obj.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType {
		if obj.Spec.UpdateStrategy.RollingUpdate == nil {
			obj.Spec.UpdateStrategy.RollingUpdate = &RollingUpdateStatefulSetStrategy{}
		}
		if obj.Spec.UpdateStrategy.RollingUpdate.Partition == nil {
			obj.Spec.UpdateStrategy.RollingUpdate.Partition = utilpointer.Int32Ptr(0)
		}
		if obj.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable == nil {
			maxUnavailable := intstr.FromInt(1)
			obj.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable = &maxUnavailable
		}
		if obj.Spec.UpdateStrategy.RollingUpdate.PodUpdatePolicy == "" {
			obj.Spec.UpdateStrategy.RollingUpdate.PodUpdatePolicy = RecreatePodUpdateStrategyType
		}
	}

	if obj.Spec.Replicas == nil {
		obj.Spec.Replicas = utilpointer.Int32Ptr(1)
	}
	if obj.Spec.RevisionHistoryLimit == nil {
		obj.Spec.RevisionHistoryLimit = utilpointer.Int32Ptr(10)
	}

	utils.SetDefaultPodTemplate(&obj.Spec.Template.Spec)
	for i := range obj.Spec.VolumeClaimTemplates {
		a := &obj.Spec.VolumeClaimTemplates[i]
		v1.SetDefaults_PersistentVolumeClaim(a)
		v1.SetDefaults_ResourceList(&a.Spec.Resources.Limits)
		v1.SetDefaults_ResourceList(&a.Spec.Resources.Requests)
		v1.SetDefaults_ResourceList(&a.Status.Capacity)
	}
}

// SetDefaults_UnitedDeployment set default values for UnitedDeployment.
func SetDefaults_UnitedDeployment(obj *UnitedDeployment) {
	if obj.Spec.Replicas == nil {
		obj.Spec.Replicas = utilpointer.Int32Ptr(1)
	}
	if obj.Spec.RevisionHistoryLimit == nil {
		obj.Spec.RevisionHistoryLimit = utilpointer.Int32Ptr(10)
	}

	if len(obj.Spec.UpdateStrategy.Type) == 0 {
		obj.Spec.UpdateStrategy.Type = ManualUpdateStrategyType
	}

	if obj.Spec.UpdateStrategy.Type == ManualUpdateStrategyType && obj.Spec.UpdateStrategy.ManualUpdate == nil {
		obj.Spec.UpdateStrategy.ManualUpdate = &ManualUpdate{}
	}

	if obj.Spec.UpdateStrategy.Type == CanaryUpdateStrategyType && obj.Spec.UpdateStrategy.CanaryUpdate == nil {
		obj.Spec.UpdateStrategy.CanaryUpdate = &CanaryUpdate{RollCount: 1, BakeTimeSeconds: 2}
	}

	if obj.Spec.Template.StatefulSetTemplate != nil {
		utils.SetDefaultPodTemplate(&obj.Spec.Template.StatefulSetTemplate.Spec.Template.Spec)
		for i := range obj.Spec.Template.StatefulSetTemplate.Spec.VolumeClaimTemplates {
			a := &obj.Spec.Template.StatefulSetTemplate.Spec.VolumeClaimTemplates[i]
			v1.SetDefaults_PersistentVolumeClaim(a)
			v1.SetDefaults_ResourceList(&a.Spec.Resources.Limits)
			v1.SetDefaults_ResourceList(&a.Spec.Resources.Requests)
			v1.SetDefaults_ResourceList(&a.Status.Capacity)
		}
	}
}

// SetDefaults_CloneSet set default values for CloneSet.
func SetDefaults_CloneSet(obj *CloneSet) {
	if obj.Spec.Replicas == nil {
		obj.Spec.Replicas = utilpointer.Int32Ptr(1)
	}
	if obj.Spec.RevisionHistoryLimit == nil {
		obj.Spec.RevisionHistoryLimit = utilpointer.Int32Ptr(10)
	}

	utils.SetDefaultPodTemplate(&obj.Spec.Template.Spec)
	for i := range obj.Spec.VolumeClaimTemplates {
		a := &obj.Spec.VolumeClaimTemplates[i]
		v1.SetDefaults_PersistentVolumeClaim(a)
		v1.SetDefaults_ResourceList(&a.Spec.Resources.Limits)
		v1.SetDefaults_ResourceList(&a.Spec.Resources.Requests)
		v1.SetDefaults_ResourceList(&a.Status.Capacity)
	}

	switch obj.Spec.UpdateStrategy.Type {
	case "":
		obj.Spec.UpdateStrategy.Type = RecreateCloneSetUpdateStrategyType
	case InPlaceIfPossibleCloneSetUpdateStrategyType, InPlaceOnlyCloneSetUpdateStrategyType:
		if obj.Spec.UpdateStrategy.InPlaceUpdateStrategy == nil {
			obj.Spec.UpdateStrategy.InPlaceUpdateStrategy = &InPlaceUpdateStrategy{}
		}
	}

	if obj.Spec.UpdateStrategy.Partition == nil {
		obj.Spec.UpdateStrategy.Partition = utilpointer.Int32Ptr(0)
	}
	if obj.Spec.UpdateStrategy.MaxUnavailable == nil {
		maxUnavailable := intstr.FromString(DefaultCloneSetMaxUnavailable)
		obj.Spec.UpdateStrategy.MaxUnavailable = &maxUnavailable
	}
	if obj.Spec.UpdateStrategy.MaxSurge == nil {
		maxSurge := intstr.FromInt(0)
		obj.Spec.UpdateStrategy.MaxSurge = &maxSurge
	}
}

// SetDefaults_DaemonSet set default values for DaemonSet.
func SetDefaults_DaemonSet(obj *DaemonSet) {
	if obj.Spec.BurstReplicas == nil {
		BurstReplicas := intstr.FromInt(250)
		obj.Spec.BurstReplicas = &BurstReplicas
	}

	if obj.Spec.UpdateStrategy.Type == "" {
		obj.Spec.UpdateStrategy.Type = RollingUpdateDaemonSetStrategyType

		// UpdateStrategy.RollingUpdate will take default values below.
		obj.Spec.UpdateStrategy.RollingUpdate = &RollingUpdateDaemonSet{}
	}

	if obj.Spec.UpdateStrategy.Type == RollingUpdateDaemonSetStrategyType {
		if obj.Spec.UpdateStrategy.RollingUpdate == nil {
			obj.Spec.UpdateStrategy.RollingUpdate = &RollingUpdateDaemonSet{}
		}
		if obj.Spec.UpdateStrategy.RollingUpdate.Partition == nil {
			obj.Spec.UpdateStrategy.RollingUpdate.Partition = new(int32)
			*obj.Spec.UpdateStrategy.RollingUpdate.Partition = 0
		}
		if obj.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable == nil {
			maxUnavailable := intstr.FromInt(1)
			obj.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable = &maxUnavailable
		}

		if obj.Spec.UpdateStrategy.RollingUpdate.Type == "" {
			obj.Spec.UpdateStrategy.RollingUpdate.Type = StandardRollingUpdateType
		}
		// Only when RollingUpdate Type is SurgingRollingUpdateType, it need to initialize the MaxSurge.
		if obj.Spec.UpdateStrategy.RollingUpdate.Type == SurgingRollingUpdateType {
			if obj.Spec.UpdateStrategy.RollingUpdate.MaxSurge == nil {
				MaxSurge := intstr.FromInt(1)
				obj.Spec.UpdateStrategy.RollingUpdate.MaxSurge = &MaxSurge
			}
		}
	}

	if obj.Spec.RevisionHistoryLimit == nil {
		obj.Spec.RevisionHistoryLimit = new(int32)
		*obj.Spec.RevisionHistoryLimit = 10
	}
}
