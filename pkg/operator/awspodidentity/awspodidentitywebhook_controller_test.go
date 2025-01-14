/*
Copyright 2018 The OpenShift Authors.

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

package awspodidentity

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	appsv1 "k8s.io/api/apps/v1"
	policyv1 "k8s.io/api/policy/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	fakeclientgo "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1 "github.com/openshift/api/config/v1"
	schemeutils "github.com/openshift/cloud-credential-operator/pkg/util"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestAWSPodIdentityWebhookController(t *testing.T) {
	schemeutils.SetupScheme(scheme.Scheme)

	getDeployment := func(c kubernetes.Interface, name, namespace string) *appsv1.Deployment {
		deployment, err := c.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err == nil {
			return deployment
		}
		return nil
	}

	getPDB := func(c kubernetes.Interface, name, namespace string) *policyv1.PodDisruptionBudget {
		pdb, err := c.PolicyV1().PodDisruptionBudgets(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err == nil {
			return pdb
		}
		return nil
	}

	tests := []struct {
		name             string
		existing         []runtime.Object
		expectErr        bool
		expectedReplicas int32
		expectPDB        bool
	}{
		{
			name: "Cluster infrastructure topology is SingleReplica",
			existing: []runtime.Object{
				&configv1.Infrastructure{
					ObjectMeta: v1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.InfrastructureStatus{
						InfrastructureTopology: configv1.SingleReplicaTopologyMode,
					},
				}},
			expectErr:        false,
			expectedReplicas: 1,
			expectPDB:        false,
		},
		{
			name: "Cluster infrastructure topology is HighlyAvailable",
			existing: []runtime.Object{
				&configv1.Infrastructure{
					ObjectMeta: v1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.InfrastructureStatus{
						InfrastructureTopology: configv1.HighlyAvailableTopologyMode,
					},
				}},
			expectErr:        false,
			expectedReplicas: 2,
			expectPDB:        true,
		},
		{
			name: "Cluster infrastructure object has no infrastructure topology set",
			existing: []runtime.Object{
				&configv1.Infrastructure{
					ObjectMeta: v1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.InfrastructureStatus{},
				}},
			expectErr:        false,
			expectedReplicas: 2,
			expectPDB:        true,
		},
		{
			name:      "Cluster infrastructure object doesn't exist",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			logger := log.WithField("controller", "awspodidentitywebhookcontrollertest")
			fakeClient := fake.NewClientBuilder().WithRuntimeObjects(test.existing...).Build()
			fakeClientset := fakeclientgo.NewSimpleClientset()
			r := &staticResourceReconciler{
				client:        fakeClient,
				clientset:     fakeClientset,
				logger:        logger,
				eventRecorder: events.NewInMemoryRecorder(""),
				cache:         resourceapply.NewResourceCache(),
				conditions:    []configv1.ClusterOperatorStatusCondition{},
				imagePullSpec: "testimagepullspec",
			}

			_, err := r.Reconcile(context.TODO(), reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "testName",
					Namespace: "testNamespace",
				},
			})

			if err != nil && !test.expectErr {
				require.NoError(t, err, "Unexpected error: %v", err)
			}
			if err == nil && test.expectErr {
				t.Errorf("Expected error but got none")
			}

			if !test.expectErr {
				podIdentityWebhookDeployment := getDeployment(fakeClientset, "pod-identity-webhook", "openshift-cloud-credential-operator")
				assert.NotNil(t, podIdentityWebhookDeployment, "did not find expected pod-identity-webhook Deployment")

				if test.expectedReplicas != 0 {
					assert.Equal(t, *podIdentityWebhookDeployment.Spec.Replicas, test.expectedReplicas, "found unexpected pod-identity-webhook deployment replicas")
				}

				podDisruptionBudget := getPDB(fakeClientset, "pod-identity-webhook", "openshift-cloud-credential-operator")
				if test.expectPDB {
					assert.NotNil(t, podDisruptionBudget, "did not find expected pod-identity-webhook PodDisruptionBudget")
				} else {
					assert.Nil(t, podDisruptionBudget, "found unexpected pod-identity-webhook PodDisruptionBudget")
				}
			}
		})
	}
}
