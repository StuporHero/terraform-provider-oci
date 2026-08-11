// Copyright (c) 2026, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package provider

import (
	"crypto/rsa"
	"testing"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeClaimConfigurationProvider stands in for the provider returned by
// oci_common_auth.OkeWorkloadIdentityConfigurationProvider, whose Region() is
// fixed to the region of the cluster the workload runs in.
type fakeClaimConfigurationProvider struct {
	region string
}

func (f fakeClaimConfigurationProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	return nil, nil
}

func (f fakeClaimConfigurationProvider) KeyID() (string, error) {
	return "fake-key-id", nil
}

func (f fakeClaimConfigurationProvider) TenancyOCID() (string, error) {
	return "fake-tenancy-ocid", nil
}

func (f fakeClaimConfigurationProvider) UserOCID() (string, error) {
	return "fake-user-ocid", nil
}

func (f fakeClaimConfigurationProvider) KeyFingerprint() (string, error) {
	return "fake-fingerprint", nil
}

func (f fakeClaimConfigurationProvider) Region() (string, error) {
	return f.region, nil
}

func (f fakeClaimConfigurationProvider) AuthType() (oci_common.AuthConfig, error) {
	return oci_common.AuthConfig{AuthType: oci_common.UnknownAuthenticationType}, nil
}

func (f fakeClaimConfigurationProvider) GetClaim(key string) (interface{}, error) {
	return "fake-claim-" + key, nil
}

func TestUnitOkeWorkloadIdentityConfigurationProviderForRegionOverridesRegion(t *testing.T) {
	delegate := fakeClaimConfigurationProvider{region: "us-ashburn-1"}

	wrapped := okeWorkloadIdentityConfigurationProviderForRegion(delegate, oci_common.StringToRegion("us-chicago-1"))

	region, err := wrapped.Region()
	require.NoError(t, err)
	assert.Equal(t, "us-chicago-1", region)
}

func TestUnitOkeWorkloadIdentityConfigurationProviderForRegionDelegatesEverythingElse(t *testing.T) {
	delegate := fakeClaimConfigurationProvider{region: "us-ashburn-1"}

	wrapped := okeWorkloadIdentityConfigurationProviderForRegion(delegate, oci_common.StringToRegion("us-chicago-1"))

	tenancy, err := wrapped.TenancyOCID()
	require.NoError(t, err)
	assert.Equal(t, "fake-tenancy-ocid", tenancy)

	keyID, err := wrapped.KeyID()
	require.NoError(t, err)
	assert.Equal(t, "fake-key-id", keyID)

	claim, err := wrapped.GetClaim("sub")
	require.NoError(t, err)
	assert.Equal(t, "fake-claim-sub", claim)
}
