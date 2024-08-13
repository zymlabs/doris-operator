// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package disaggregated_fe

import (
	dv1 "github.com/selectdb/doris-operator/api/disaggregated/cluster/v1"
	"github.com/selectdb/doris-operator/pkg/common/utils/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (dfc *DisaggregatedFEController) newHeadlessService(ddc *dv1.DorisDisaggregatedCluster, cvs map[string]interface{}) *corev1.Service {
	ddcService := ddc.Spec.FeSpec.CommonSpec.Service
	ports := newFEServicePorts(cvs, ddcService)
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ddc.GetFEServiceName() + "-headless",
			Namespace: ddc.Namespace,
			Labels:    dfc.newFESchedulerLabels(ddc.Namespace),
		},
		Spec: corev1.ServiceSpec{
			Selector:        dfc.newFEPodsSelector(ddc.Name),
			Type:            corev1.ServiceTypeClusterIP,
			ClusterIP:       "None",
			Ports:           ports,
			SessionAffinity: corev1.ServiceAffinityClientIP,
		},
	}

	if ddcService != nil && ddcService.Type != "" {
		svc.Spec.Type = ddcService.Type
	}
	if ddcService != nil {
		svc.Annotations = ddcService.Annotations
	}

	// The external load balancer provided by the cloud provider may cause the client IP received by the service to change.
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		svc.Spec.SessionAffinity = corev1.ServiceAffinityNone
	}
	svc.OwnerReferences = []metav1.OwnerReference{resource.GetOwnerReference(ddc)}

	return &svc
}
