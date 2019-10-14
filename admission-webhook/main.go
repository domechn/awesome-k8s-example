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
	"flag"
	"os"

	"k8s.io/klog"
	cfg "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	certDir = "/etc/admission/cert"
	port    = 8080
)

func init() {
	flag.StringVar(&certDir, "cert-dir", certDir, "File containing the default x509 Certificate and x509 private key for HTTPS.")
	flag.IntVar(&port, "port", port, "Listen port of webhook server")
}

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()
	flag.Parse()

	c := cfg.GetConfigOrDie()
	managerOptions := manager.Options{
		Port:    port,
		CertDir: certDir,
	}
	mgr, err := manager.New(c, managerOptions)
	if err != nil {
		klog.Fatal("unable to set up overall controller manager", err)
	}

	hookServer := mgr.GetWebhookServer()
	hookServer.Register("/mutate", &webhook.Admission{Handler: &simpleMutating{}})
	hookServer.Register("/validate", &webhook.Admission{Handler: &simpleValidating{}})

	klog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Error(err, "unable to run manager", err)
		os.Exit(1)
	}

}
