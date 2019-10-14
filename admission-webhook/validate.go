/*
Copyright (c) 2019 Domgoer Inc.
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

package k8s_admission_webhook_example

import (
	"context"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type simpleValidating struct {
	client  client.Client
	decoder *admission.Decoder
}

var _ admission.Handler = &simpleValidating{}

func (s *simpleValidating) Handle(ctx context.Context, req admission.Request) admission.Response {
	obj := &corev1.Pod{}

	err := s.decoder.Decode(req, obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	allowed, reason, err := validate(ctx, obj)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

func validate(ctx context.Context, obj *corev1.Pod) (bool, string, error) {
	panic(" do some validating")
}

// simpleValidating implements inject.Client.
// A client will be automatically injected.

// InjectClient injects the client.
func (s *simpleValidating) InjectClient(c client.Client) error {
	s.client = c
	return nil
}

// simpleValidating implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (s *simpleValidating) InjectDecoder(d *admission.Decoder) error {
	s.decoder = d
	return nil
}
