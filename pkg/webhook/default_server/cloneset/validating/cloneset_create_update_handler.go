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

package validating

import (
	"context"
	"net/http"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validating-create-update-cloneset"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &CloneSetCreateUpdateHandler{})
}

// CloneSetCreateUpdateHandler handles CloneSet
type CloneSetCreateUpdateHandler struct {
	Client client.Client

	// Decoder decodes objects
	Decoder types.Decoder
}

var _ admission.Handler = &CloneSetCreateUpdateHandler{}

// Handle handles admission requests.
func (h *CloneSetCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &appsv1alpha1.CloneSet{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	switch req.AdmissionRequest.Operation {
	case admissionv1beta1.Create:
		if allErrs := h.validateCloneSet(obj); len(allErrs) > 0 {
			return admission.ErrorResponse(http.StatusUnprocessableEntity, allErrs.ToAggregate())
		}
	case admissionv1beta1.Update:
		oldObj := &appsv1alpha1.CloneSet{}
		if err := h.Decoder.Decode(types.Request{AdmissionRequest: &admissionv1beta1.AdmissionRequest{Object: req.AdmissionRequest.OldObject}}, oldObj); err != nil {
			return admission.ErrorResponse(http.StatusBadRequest, err)
		}
		if allErrs := h.validateCloneSetUpdate(obj, oldObj); len(allErrs) > 0 {
			return admission.ErrorResponse(http.StatusUnprocessableEntity, allErrs.ToAggregate())
		}
	}

	return admission.ValidationResponse(true, "")
}

var _ inject.Client = &CloneSetCreateUpdateHandler{}

// InjectClient injects the client into the CloneSetCreateUpdateHandler
func (h *CloneSetCreateUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &CloneSetCreateUpdateHandler{}

// InjectDecoder injects the decoder into the CloneSetCreateUpdateHandler
func (h *CloneSetCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
