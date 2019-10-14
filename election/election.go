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

package election

import (
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

const (
	LockType = "configmaps"

	LeaseDuration = time.Second * 30
	RetryPeriod   = time.Second * 2
	RenewDeadline = time.Second * 10
)

type Config struct {
	ResourceName      string
	ResourceNamespace string

	Callbacks leaderelection.LeaderCallbacks
}

// NewElection returns leaderelection.LeaderElector to start election, should use leaderelection.LeaderElector.Run(ctx)
func NewElection(config Config, client kubernetes.Interface) (*leaderelection.LeaderElector, error) {
	lec, err := getLeaderElectionConfig(config, client)
	if err != nil {
		return nil, err
	}
	return leaderelection.NewLeaderElector(lec)
}

func getLeaderElectionConfig(config Config, client kubernetes.Interface) (lec leaderelection.LeaderElectionConfig, err error) {
	leaderElectionBroadcaster := record.NewBroadcaster()
	host, err := os.Hostname()
	if err != nil {
		return
	}
	recorder := leaderElectionBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: config.ResourceName, Host: host})
	id := string(uuid.NewUUID())

	rl, err := resourcelock.New(LockType,
		config.ResourceNamespace,
		config.ResourceName,
		client.CoreV1(),
		client.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: recorder,
		})
	if err != nil {
		return
	}

	lec = leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: LeaseDuration,
		RenewDeadline: RenewDeadline,
		RetryPeriod:   RetryPeriod,
		Callbacks:     config.Callbacks,
		WatchDog:      leaderelection.NewLeaderHealthzAdaptor(time.Second * 20),
		Name:          config.ResourceName,
	}
	return
}
