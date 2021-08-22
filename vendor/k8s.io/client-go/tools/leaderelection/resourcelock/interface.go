/*
Copyright 2016 The Kubernetes Authors.

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

package resourcelock

import (
	"context"
	"fmt"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	LeaderElectionRecordAnnotationKey = "control-plane.alpha.kubernetes.io/leader"
	EndpointsResourceLock             = "endpoints"
	ConfigMapsResourceLock            = "configmaps"
	LeasesResourceLock                = "leases"
	EndpointsLeasesResourceLock       = "endpointsleases"
	ConfigMapsLeasesResourceLock      = "configmapsleases"
)

// LeaderElectionRecord is the record that is stored in the leader election annotation.
// This information should be used for observational purposes only and could be replaced
// with a random string (e.g. UUID) with only slight modification of this code.
// TODO(mikedanese): this should potentially be versioned
// key (分布式锁)存在于Etcd 集群的/registry/services/endpoints/kube-system/kube-scheduler中。
// 该key中存储的是竞选为领导者节点的信息，它通过LeaderElectionRecord结构体进行描述:
type LeaderElectionRecord struct {
	// HolderIdentity is the ID that owns the lease. If empty, no one owns this lease and
	// all callers may acquire. Versions of this library prior to Kubernetes 1.14 will not
	// attempt to acquire leases with empty identities and will wait for the full lease
	// interval to expire before attempting to reacquire. This value is set to empty when
	// a client voluntarily steps down.
	// 领导者身份标识，通常为Hostname_<hash 值>.
	HolderIdentity string `json:"holderIdentity"`
	// 领导者租约的时长。
	LeaseDurationSeconds int `json:"leaseDurationSeconds"`
	// 领导者获得锁的时间。
	AcquireTime metav1.Time `json:"acquireTime"`
	// 领导者续租的时间。
	RenewTime metav1.Time `json:"renewTime"`
	// 领导者选举切换的次数。
	LeaderTransitions int `json:"leaderTransitions"`
}

// EventRecorder records a change in the ResourceLock.
type EventRecorder interface {
	Eventf(obj runtime.Object, eventType, reason, message string, args ...interface{})
}

// ResourceLockConfig common data that exists across different
// resource locks
type ResourceLockConfig struct {
	// Identity is the unique string identifying a lease holder across
	// all participants in an election.
	Identity string
	// EventRecorder is optional.
	EventRecorder EventRecorder
}

// Interface offers a common interface for locking on arbitrary
// resources used in leader election.  The Interface is used
// to hide the details on specific implementations in order to allow
// them to change over time.  This interface is strictly for use
// by the leaderelection code.
// 每种资源锁实现了对key (资源锁)的操作方法，它的接口定义如下:
type Interface interface {
	// Get returns the LeaderElectionRecord
	// 用于获取资源锁的所有信息
	Get(ctx context.Context) (*LeaderElectionRecord, []byte, error)

	// Create attempts to create a LeaderElectionRecord
	// 用于创建资源锁
	Create(ctx context.Context, ler LeaderElectionRecord) error

	// Update will update and existing LeaderElectionRecord
	// 用于更新资源锁信息
	Update(ctx context.Context, ler LeaderElectionRecord) error

	// RecordEvent is used to record events
	// 通过EventBroadcaster事件管理器记录事件
	RecordEvent(string)

	// Identity will return the locks Identity
	// 用于获取领导者身份标识
	Identity() string

	// Describe is used to convert details on current resource lock
	// into a string
	// 用于获取资源锁的信息。
	Describe() string
}

// Manufacture will create a lock of a given type according to the input parameters
func New(lockType string, ns string, name string, coreClient corev1.CoreV1Interface, coordinationClient coordinationv1.CoordinationV1Interface, rlc ResourceLockConfig) (Interface, error) {
	endpointsLock := &EndpointsLock{
		EndpointsMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
		Client:     coreClient,
		LockConfig: rlc,
	}
	configmapLock := &ConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
		Client:     coreClient,
		LockConfig: rlc,
	}
	leaseLock := &LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
		Client:     coordinationClient,
		LockConfig: rlc,
	}
	switch lockType {
	case EndpointsResourceLock:
		return endpointsLock, nil
	case ConfigMapsResourceLock:
		return configmapLock, nil
	case LeasesResourceLock:
		return leaseLock, nil
	case EndpointsLeasesResourceLock:
		return &MultiLock{
			Primary:   endpointsLock,
			Secondary: leaseLock,
		}, nil
	case ConfigMapsLeasesResourceLock:
		return &MultiLock{
			Primary:   configmapLock,
			Secondary: leaseLock,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid lock-type %s", lockType)
	}
}

// NewFromKubeconfig will create a lock of a given type according to the input parameters.
// Timeout set for a client used to contact to Kubernetes should be lower than
// RenewDeadline to keep a single hung request from forcing a leader loss.
// Setting it to max(time.Second, RenewDeadline/2) as a reasonable heuristic.
func NewFromKubeconfig(lockType string, ns string, name string, rlc ResourceLockConfig, kubeconfig *restclient.Config, renewDeadline time.Duration) (Interface, error) {
	// shallow copy, do not modify the kubeconfig
	config := *kubeconfig
	timeout := renewDeadline / 2
	if timeout < time.Second {
		timeout = time.Second
	}
	config.Timeout = timeout
	leaderElectionClient := clientset.NewForConfigOrDie(restclient.AddUserAgent(&config, "leader-election"))
	return New(lockType, ns, name, leaderElectionClient.CoreV1(), leaderElectionClient.CoordinationV1(), rlc)
}
