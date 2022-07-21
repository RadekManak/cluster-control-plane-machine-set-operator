/*
Copyright 2022 Red Hat, Inc.

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

package resourcebuilder

import (
	"encoding/json"

	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func OpenStackProviderSpec() OpenStackProviderSpecBuilder {
	return OpenStackProviderSpecBuilder{
		AvailabilityZone: "TODO",
	}
}

type OpenStackProviderSpecBuilder struct {
	AvailabilityZone string
}

func (o OpenStackProviderSpecBuilder) Build() *machinev1alpha1.OpenstackProviderSpec {
	// TODO
	return &machinev1alpha1.OpenstackProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "openstackproviderconfig.openshift.io/v1alpha1",
			Kind:       "OpenstackProviderSpec",
		},
		CloudName: "openstack",
		Flavor:    "ci.m1.xlarge",
		Image:     "0bnhphb-b5564-2wmsh-rhcos",
		Networks: []machinev1alpha1.NetworkParam{
			{
				Subnets: []machinev1alpha1.SubnetParam{
					{
						Filter: machinev1alpha1.SubnetFilter{
							Name: "0bnhphb-b5564-2wmsh-nodes",
							Tags: "openshiftClusterID=0bnhphb-b5564-2wmsh",
						},
					},
				},
			},
		},
		AvailabilityZone: o.AvailabilityZone,
		SecurityGroups: []machinev1alpha1.SecurityGroupParam{
			{
				Name: " 0bnhphb-b5564-2wmsh-worker",
			},
		},
		UserDataSecret: &corev1.SecretReference{
			Name: "worker-user-data",
		},
		Trunk: true,
		Tags: []string{
			"openshiftClusterID=0bnhphb-b5564-2wmsh",
		},
		ServerMetadata: map[string]string{
			"Name":               "0bnhphb-b5564-2wmsh-worker",
			"openshiftClusterID": "0bnhphb-b5564-2wmsh",
		},
		ServerGroupName: "0bnhphb-b5564-2wmsh-worker",
	}
}

// BuildRawExtension builds a new OpenStack machine config based on the configuration of the builder.
func (o OpenStackProviderSpecBuilder) BuildRawExtension() *runtime.RawExtension {
	providerConfig := o.Build()

	raw, err := json.Marshal(providerConfig)
	if err != nil {
		// As we are building the input to json.Marshal, this should never happen.
		panic(err)
	}

	return &runtime.RawExtension{
		Raw: raw,
	}
}

// WithAvailabilityZone sets the availability zone for the OpenStack machine config builder.
func (o OpenStackProviderSpecBuilder) WithAvailabilityZone(az string) OpenStackProviderSpecBuilder {
	o.AvailabilityZone = az
	return o
}
