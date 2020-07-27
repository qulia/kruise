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

package statefulsetrollout

import (
	"testing"

	"github.com/onsi/gomega"
	appsv1alpha1 "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}

func TestReconcile(t *testing.T) {
	g, requests, stopMgr, mgrStopped := setUp(t)
	defer func() {
		clean(g, c)
		close(stopMgr)
		mgrStopped.Wait()
	}()

	instance := &appsv1alpha1.StatefulSetRollout{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: appsv1alpha1.StatefulSetRolloutSpec{
			Replicas: &one,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "foo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "foo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container-a",
							Image: "nginx:1.0",
						},
					},
				},
			},
			Rollout: appsv1alpha1.StatefulSetRolloutStrategy{Type: "blue-green"},
		},
	}
	// Create the UnitedDeployment object and expect the Reconcile and Deployment to be created
	err := c.Create(context.TODO(), instance)
	defer c.Delete(context.TODO(), instance)

	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
}
