/*
Copyright 2019 The Crossplane Authors.

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

package redis

import (
	"testing"

	redismgmt "github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	azure "github.com/crossplaneio/stack-azure/pkg/clients"
	"github.com/google/go-cmp/cmp"

	"github.com/crossplaneio/stack-azure/apis/cache/v1alpha3"
)

const (
	skuName     = "basic"
	skuFamily   = "C"
	skuCapacity = 1
)

var (
	location           = "us-east1"
	zones              = []string{"us-east1a", "us-east1b"}
	tags               = map[string]string{"key1": "val1"}
	tags2              = map[string]string{"key1": "val1", "key2": "val2"}
	enableNonSSLPort   = true
	subnetID           = "coolsubnet"
	staticIP           = "172.16.0.1"
	shardCount         = 3
	redisConfiguration = map[string]string{"cool": "socool"}
	tenantSettings     = map[string]string{"tenant1": "is-crazy"}
	tenantSettings2    = map[string]string{"tenant1": "is-crazy", "tenant2": "not-so-much"}
	minTLSVersion      = "1.1"

	redisVersion  = "3.2"
	hostName      = "108.8.8.1"
	port          = 6374
	sslPort       = 453
	linkedServers = []string{"server1", "server2"}
)

func TestNewCreateParameters(t *testing.T) {
	cases := []struct {
		name string
		r    *v1alpha3.Redis
		want redismgmt.CreateParameters
	}{
		{
			name: "Successful",
			r: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						Location: location,
						Zones:    zones,
						Tags:     tags,
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						SubnetID:           &subnetID,
						StaticIP:           &staticIP,
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						TenantSettings:     tenantSettings,
						ShardCount:         &shardCount,
						MinimumTLSVersion:  &minTLSVersion,
					},
				},
			},
			want: redismgmt.CreateParameters{
				Location: azure.ToStringPtr(location),
				Zones:    azure.ToStringArrayPtr(zones),
				Tags:     azure.ToStringPtrMap(tags),
				CreateProperties: &redismgmt.CreateProperties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					SubnetID:           azure.ToStringPtr(subnetID),
					StaticIP:           azure.ToStringPtr(staticIP),
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					TenantSettings:     azure.ToStringPtrMap(tenantSettings),
					ShardCount:         azure.ToInt32Ptr(shardCount),
					MinimumTLSVersion:  redismgmt.TLSVersion(minTLSVersion),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewCreateParameters(tc.r)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NewCreateParameters(...): -want, +got\n%s", diff)
			}
		})
	}
}

func TestNewUpdateParameters(t *testing.T) {
	cases := []struct {
		name string
		r    *v1alpha3.Redis
		want redismgmt.UpdateParameters
	}{
		{
			name: "UpdatableFieldsOnly",
			r: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						Tags: tags,
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
						TenantSettings:     tenantSettings,
						MinimumTLSVersion:  &minTLSVersion,
					},
				},
			},
			want: redismgmt.UpdateParameters{
				Tags: azure.ToStringPtrMap(tags),
				UpdateProperties: &redismgmt.UpdateProperties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount),
					TenantSettings:     azure.ToStringPtrMap(tenantSettings),
					MinimumTLSVersion:  redismgmt.TLSVersion(minTLSVersion),
				},
			},
		},
		{
			name: "SuperfluousFields",
			r: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						SubnetID:           &subnetID,
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,

						// StaticIP cannot be updated and should be omitted.
						StaticIP: &staticIP,
					},
				},
			},
			want: redismgmt.UpdateParameters{
				UpdateProperties: &redismgmt.UpdateProperties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewUpdateParameters(tc.r)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NewUpdateParameters(...): -want, +got\n%s", diff)
			}
		})
	}
}

func TestNeedsUpdate(t *testing.T) {
	cases := []struct {
		name string
		kube *v1alpha3.Redis
		az   redismgmt.ResourceType
		want bool
	}{
		{
			name: "NeedsLessCapacity",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
					},
				},
			},
			az: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity + 1),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount),
				},
			},
			want: true,
		},
		{
			name: "NeedsNewRedisConfiguration",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
					},
				},
			},
			az: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(map[string]string{"super": "cool"}),
					ShardCount:         azure.ToInt32Ptr(shardCount),
				},
			},
			want: true,
		},
		{
			name: "NeedsSSLPortDisabled",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
					},
				},
			},
			az: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(!enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount),
				},
			},
			want: true,
		},
		{
			name: "NeedsFewerShards",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
					},
				},
			},
			az: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount + 1),
				},
			},
			want: true,
		},
		{
			name: "NeedsDifferentTenantSettings",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
						TenantSettings:     tenantSettings2,
					},
				},
			},
			az: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount + 1),
					TenantSettings:     azure.ToStringPtrMap(tenantSettings),
				},
			},
			want: true,
		},
		{
			name: "NeedsTagsToBeUpdated",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
						Tags:               tags2,
					},
				},
			},
			az: redismgmt.ResourceType{
				Tags: azure.ToStringPtrMap(tags),
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount + 1),
				},
			},
			want: true,
		},
		{
			name: "NeedsNoUpdate",
			kube: &v1alpha3.Redis{
				Spec: v1alpha3.RedisSpec{
					ForProvider: v1alpha3.RedisParameters{
						SKU: v1alpha3.SKU{
							Name:     skuName,
							Family:   skuFamily,
							Capacity: skuCapacity,
						},
						EnableNonSSLPort:   &enableNonSSLPort,
						RedisConfiguration: redisConfiguration,
						ShardCount:         &shardCount,
					},
				},
			},
			az: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount),
				},
			},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NeedsUpdate(tc.kube, tc.az)
			if got != tc.want {
				t.Errorf("NeedsUpdate(...): want %t, got %t", tc.want, got)
			}
		})
	}
}

func TestGenerateObservation(t *testing.T) {
	cases := map[string]struct {
		arg  redismgmt.ResourceType
		want v1alpha3.RedisObservation
	}{
		"FullConversion": {
			arg: redismgmt.ResourceType{
				Properties: &redismgmt.Properties{
					RedisVersion:      azure.ToStringPtr(redisVersion),
					ProvisioningState: redismgmt.ProvisioningState(ProvisioningStateCreating),
					HostName:          azure.ToStringPtr(hostName),
					Port:              azure.ToInt32(&port),
					LinkedServers: &[]redismgmt.LinkedServer{
						{ID: azure.ToStringPtr(linkedServers[0])},
						{ID: azure.ToStringPtr(linkedServers[1])},
					},
					Sku: &redismgmt.Sku{
						Name:     redismgmt.SkuName(skuName),
						Family:   redismgmt.SkuFamily(skuFamily),
						Capacity: azure.ToInt32Ptr(skuCapacity),
					},
					SubnetID:           azure.ToStringPtr(subnetID),
					StaticIP:           azure.ToStringPtr(staticIP),
					TenantSettings:     azure.ToStringPtrMap(tenantSettings),
					MinimumTLSVersion:  redismgmt.TLSVersion(minTLSVersion),
					EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
					RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
					ShardCount:         azure.ToInt32Ptr(shardCount),
					SslPort:            azure.ToInt32(&sslPort),
				},
			},
			want: v1alpha3.RedisObservation{
				RedisVersion:       redisVersion,
				ProvisioningState:  ProvisioningStateCreating,
				HostName:           hostName,
				Port:               port,
				SSLPort:            sslPort,
				RedisConfiguration: redisConfiguration,
				EnableNonSSLPort:   enableNonSSLPort,
				TenantSettings:     tenantSettings,
				ShardCount:         shardCount,
				MinimumTLSVersion:  minTLSVersion,
				LinkedServers:      linkedServers,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateObservation(tc.arg)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateObservation(...): -want, +got\n%s", diff)
			}
		})
	}
}

func TestLateInitialize(t *testing.T) {
	type args struct {
		az   redismgmt.ResourceType
		spec *v1alpha3.RedisParameters
	}
	type want struct {
		spec *v1alpha3.RedisParameters
	}
	cases := map[string]struct {
		args
		want
	}{
		"LateInitializeEmptyObject": {
			args: args{
				az: redismgmt.ResourceType{
					Zones: azure.ToStringArrayPtr(zones),
					Tags:  azure.ToStringPtrMap(tags),
					Properties: &redismgmt.Properties{
						RedisVersion:      azure.ToStringPtr(redisVersion),
						ProvisioningState: redismgmt.ProvisioningState(ProvisioningStateCreating),
						HostName:          azure.ToStringPtr(hostName),
						Port:              azure.ToInt32(&port),
						LinkedServers: &[]redismgmt.LinkedServer{
							{ID: azure.ToStringPtr(linkedServers[0])},
							{ID: azure.ToStringPtr(linkedServers[1])},
						},
						Sku: &redismgmt.Sku{
							Name:     redismgmt.SkuName(skuName),
							Family:   redismgmt.SkuFamily(skuFamily),
							Capacity: azure.ToInt32Ptr(skuCapacity),
						},
						SubnetID:           azure.ToStringPtr(subnetID),
						StaticIP:           azure.ToStringPtr(staticIP),
						TenantSettings:     azure.ToStringPtrMap(tenantSettings),
						MinimumTLSVersion:  redismgmt.TLSVersion(minTLSVersion),
						EnableNonSslPort:   azure.ToBoolPtr(enableNonSSLPort),
						RedisConfiguration: azure.ToStringPtrMap(redisConfiguration),
						ShardCount:         azure.ToInt32Ptr(shardCount),
						SslPort:            azure.ToInt32(&sslPort),
					},
				},
				spec: &v1alpha3.RedisParameters{},
			},
			want: want{
				spec: &v1alpha3.RedisParameters{
					Zones:              zones,
					Tags:               tags,
					SubnetID:           &subnetID,
					StaticIP:           &staticIP,
					EnableNonSSLPort:   &enableNonSSLPort,
					RedisConfiguration: redisConfiguration,
					TenantSettings:     tenantSettings,
					ShardCount:         &shardCount,
					MinimumTLSVersion:  &minTLSVersion,
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			LateInitialize(tc.args.spec, tc.args.az)
			if diff := cmp.Diff(tc.want.spec, tc.args.spec); diff != "" {
				t.Errorf("LateInitialize(...): -want, +got\n%s", diff)
			}
		})
	}
}
