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

package providerconfig

import (
	"encoding/json"
	"fmt"

	v1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/api/machine/v1alpha1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// OpenStackProviderConfig holds the provider spec of an OpenStack Machine.
// It allows external code to extract and inject failure domain information,
// as well as gathering the stored config.
type OpenStackProviderConfig struct {
	providerConfig machinev1alpha1.OpenstackProviderSpec
}

// InjectFailureDomain return a new OpenStackProviderConfig configured with the provided failure domain
// information.
func (o OpenStackProviderConfig) InjectFailureDomain(fd machinev1.OpenStackFailureDomain) OpenStackProviderConfig {
	newOpenstackProviderConfig := o
	newOpenstackProviderConfig.providerConfig.AvailabilityZone = fd.ComputeZone
	newOpenstackProviderConfig.providerConfig.RootVolume.Zone = fd.StorageZone

	if fd.Subnet.UUID != "" {
		_, err := findSubnetByUUID(o.providerConfig.Networks, fd.Subnet.UUID)
		if err == nil {
			newOpenstackProviderConfig.providerConfig.Networks = append(newOpenstackProviderConfig.providerConfig.Networks, machinev1alpha1.NetworkParam{
				Filter: machinev1alpha1.Filter{},
				Subnets: []machinev1alpha1.SubnetParam{
					{
						UUID: fd.Subnet.UUID,
					},
				},
			})
		}
	}

	if fd.Subnet.Filter.Name != "" {
		_, err := findSubnetByName(o.providerConfig.Networks, fd.Subnet.Filter.Name)
		if err == nil {
			newOpenstackProviderConfig.providerConfig.Networks = append(newOpenstackProviderConfig.providerConfig.Networks, machinev1alpha1.NetworkParam{
				Filter: machinev1alpha1.Filter{},
				Subnets: []machinev1alpha1.SubnetParam{
					{
						Filter: fd.Subnet.Filter,
					},
				},
			})
		}
	}

	return newOpenstackProviderConfig
}

func findSubnetByUUID(networks []machinev1alpha1.NetworkParam, subnetUUID string) (v1alpha1.SubnetParam, error) {
	for _, network := range networks {
		for _, subnet := range network.Subnets {
			if subnet.UUID == subnetUUID {
				return subnet, nil
			}
		}
	}
	return machinev1alpha1.SubnetParam{}, fmt.Errorf("Primary subnet %s not specified on machine", subnetUUID)
}

func findSubnetByName(networks []machinev1alpha1.NetworkParam, subnetName string) (v1alpha1.SubnetParam, error) {
	for _, network := range networks {
		for _, subnet := range network.Subnets {
			if subnet.Filter.Name == subnetName {
				return subnet, nil
			}
		}
	}
	return v1alpha1.SubnetParam{}, fmt.Errorf("Primary subnet %s not specified on machine", subnetName)
}

// ExtractFailureDomain returns an OpenStackFailureDomain based on the failure domain
// information stored within the OpenStackProviderConfig.
func (o OpenStackProviderConfig) ExtractFailureDomain() machinev1.OpenStackFailureDomain {
	fd := machinev1.OpenStackFailureDomain{
		ComputeZone: o.providerConfig.AvailabilityZone,
		StorageZone: o.providerConfig.RootVolume.Zone,
	}

	if o.providerConfig.PrimarySubnet != "" {
		fd.Subnet = machinev1alpha1.SubnetParam{UUID: o.providerConfig.PrimarySubnet}
	} else {
		if len(o.providerConfig.Networks) > 0 && len(o.providerConfig.Networks[0].Subnets) > 0 {
			fd.Subnet = o.providerConfig.Networks[0].Subnets[0]
		}
	}
	return fd
}

// Config returns the stored OpenstackProviderSpec.
func (o OpenStackProviderConfig) Config() machinev1alpha1.OpenstackProviderSpec {
	return o.providerConfig
}

// newOpenStackProviderConfig creates an OpenStack type ProviderConfig fom the raw extension.
// It should return an error if the provided RawExtension does not represent
// an OpenstackProviderSpec.
func newOpenStackProviderConfig(raw *runtime.RawExtension) (ProviderConfig, error) {
	openstackProviderSpec := machinev1alpha1.OpenstackProviderSpec{}
	if err := json.Unmarshal(raw.Raw, &openstackProviderSpec); err != nil {
		return nil, fmt.Errorf("could not unmarshal provider spec: %w", err)
	}

	openstackProviderConfig := OpenStackProviderConfig{
		providerConfig: openstackProviderSpec,
	}

	config := providerConfig{
		platformType: v1.OpenStackPlatformType,
		openstack:    openstackProviderConfig,
	}

	return config, nil
}
