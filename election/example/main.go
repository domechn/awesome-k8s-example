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

package main

import (
	"context"
	"fmt"
	"github.com/domgoer/awesome-k8s-example/election-example"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
)

func main() {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", "/Users/dmc/.kube/config")
		if err != nil {
			panic(err)
		}
	}
	client := kubernetes.NewForConfigOrDie(cfg)
	le, err := election.NewElection(election.Config{
		ResourceName:      "election-example",
		ResourceNamespace: "default",
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(i context.Context) {
				fmt.Println("start leader election")
			},
			OnStoppedLeading: func() {
				fmt.Println("stop leader election")
			},
			OnNewLeader: func(id string) {
				fmt.Println("new leader: ", id)
			},
		},
	}, client)
	if err != nil {
		panic(err)
	}

	le.Run(context.Background())
	select {}
}
