// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package store

import (
	appsv1 "k8s.io/client-go/listers/apps/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
	v1beta1 "k8s.io/client-go/listers/extensions/v1beta1"
)

//Lister kube-api client cache
type Lister struct {
	Ingress     v1beta1.IngressLister
	Service     corev1.ServiceLister
	Secret      corev1.SecretLister
	StatefulSet appsv1.StatefulSetLister
	Deployment  appsv1.DeploymentLister
	Pod         corev1.PodLister
	ConfigMap   corev1.ConfigMapLister
}
