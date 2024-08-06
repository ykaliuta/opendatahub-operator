/*
Copyright 2023.

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

package main

import (
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/opendatahub-io/opendatahub-operator/v2/controllers/webhook"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		WebhookServer: ctrlwebhook.NewServer(ctrlwebhook.Options{
			Port: 9443,
			// TLSOpts: , // TODO: it was not set in the old code
		}),
	})
	if err != nil {
		fmt.Printf("unable to setup manager: %v\n", err)
		os.Exit(1)
	}

	(&webhook.OpenDataHubValidatingWebhook{
		Client:  mgr.GetClient(),
		Decoder: admission.NewDecoder(mgr.GetScheme()),
	}).SetupWithManager(mgr)

	(&webhook.OpenDataHubMutatingWebhook{
		Client:  mgr.GetClient(),
		Decoder: admission.NewDecoder(mgr.GetScheme()),
	}).SetupWithManager(mgr)

	fmt.Println("Starting webhook server")

	if err := mgr.Start(ctx); err != nil {
		fmt.Printf("unable to start manager: %v\n", err)
		os.Exit(1)
	}
}
