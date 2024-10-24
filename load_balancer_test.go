package elestio

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupLoadBalancerTestCase(t *testing.T) *Client {
	t.Log("Setup service test case")

	email, apiKey := os.Getenv("ELESTIO_INTEGRATION_EMAIL"), os.Getenv("ELESTIO_INTEGRATION_API_KEY")

	c, err := NewClient(email, apiKey)
	if err != nil {
		t.Fatalf("client authentication error: %s", err)
	}

	return c
}

func TestLoadBalancerHandler_Get(t *testing.T) {
	t.Skip("Skipping test")
	c := setupLoadBalancerTestCase(t)

	projectID := "596"
	loadBalancerID := "54654320"

	loadBalancer, err := c.LoadBalancer.Get(projectID, loadBalancerID)
	require.NoError(t, err, "expected no error when getting loadBalancer")
	require.NotNil(t, loadBalancer, "expected non-nil loadBalancer")
	require.Equal(t, loadBalancerID, loadBalancer.ID, "expected loadBalancer ID to be "+loadBalancerID)

	fmt.Fprintf(os.Stdout, "LoadBalancer: %v", loadBalancer)
}

func TestLoadBalancerHandler_Create(t *testing.T) {
	t.Skip("Skipping test")
	c := setupLoadBalancerTestCase(t)

	projectId := "596"

	loadBalancer, err := c.LoadBalancer.Create(CreateLoadBalancerRequest{
		ProjectID:    projectId,
		ProviderName: "hetzner",
		Datacenter:   "fsn1",
		ServerType:   "MEDIUM-2C-4G",
		Config: CreateLoadBalancerRequestConfig{
			HostHeader:             "$http_host",
			IsAccessLogsEnabled:    true,
			IsForceHTTPSEnabled:    true,
			IPRateLimit:            100,
			IsIPRateLimitEnabled:   false,
			OutputCacheInSeconds:   0,
			IsStickySessionEnabled: false,
			IsProxyProtocolEnabled: false,
			SSLDomains:             []string{},
			ForwardRules: []LoadBalancerConfigForwardRule{
				{
					Protocol:       "HTTP",
					TargetProtocol: "HTTP",
					Port:           "80",
					TargetPort:     "3000",
				},
				{
					Protocol:       "HTTPS",
					TargetProtocol: "HTTP",
					Port:           "443",
					TargetPort:     "3000",
				},
			},
			OutputHeaders:         []LoadBalancerConfigOutputHeader{},
			TargetServices:        []string{"elest.io"},
			RemoveResponseHeaders: []string{},
		},
	})

	require.NoError(t, err, "expected no error when creating loadBalancer")
	require.NotNil(t, loadBalancer, "expected non-nil loadBalancer")
	fmt.Fprintf(os.Stdout, "LoadBalancer: %v", loadBalancer)
	require.Equal(t, "MEDIUM-2C-4G", loadBalancer.ServerType, "expected loadBalancer server type to be MEDIUM-2C-4G")
}

func TestLoadBalancerHandler_UpdateConfig(t *testing.T) {
	t.Skip("Skipping test")
	c := setupLoadBalancerTestCase(t)

	projectID := "596"
	loadBalancerID := "54654320"

	loadBalancer, err := c.LoadBalancer.Get(projectID, loadBalancerID)
	require.NoError(t, err, "expected no error when getting initial loadBalancer")

	updatedLoadBalancer, err := c.LoadBalancer.UpdateConfig(projectID, loadBalancer.ID, UpdateLoadBalancerConfigRequest{
		HostHeader:             loadBalancer.Config.HostHeader,
		IsAccessLogsEnabled:    !loadBalancer.Config.IsAccessLogsEnabled,
		IsForceHTTPSEnabled:    loadBalancer.Config.IsForceHTTPSEnabled,
		IPRateLimit:            loadBalancer.Config.IPRateLimit,
		IsIPRateLimitEnabled:   loadBalancer.Config.IsIPRateLimitEnabled,
		OutputCacheInSeconds:   loadBalancer.Config.OutputCacheInSeconds,
		IsStickySessionEnabled: loadBalancer.Config.IsStickySessionEnabled,
		IsProxyProtocolEnabled: loadBalancer.Config.IsProxyProtocolEnabled,
		SSLDomains:             loadBalancer.Config.SSLDomains,
		ForwardRules:           loadBalancer.Config.ForwardRules,
		OutputHeaders:          loadBalancer.Config.OutputHeaders,
		TargetServices:         loadBalancer.Config.TargetServices,
		RemoveResponseHeaders:  loadBalancer.Config.RemoveResponseHeaders,
	})

	require.NoError(t, err, "expected no error when updating loadBalancer")
	require.NotNil(t, loadBalancer, "expected non-nil loadBalancer")
	require.Equal(t, !loadBalancer.Config.IsAccessLogsEnabled, updatedLoadBalancer.Config.IsAccessLogsEnabled, "expected loadBalancer config IsAccessLogsEnabled to be updated")
}

func TestLoadBalancerHandler_Delete(t *testing.T) {
	t.Skip("Skipping test")
	c := setupLoadBalancerTestCase(t)

	projectID := "596"
	loadBalancerID := "54654320"

	err := c.LoadBalancer.Delete(projectID, loadBalancerID, false)
	require.NoError(t, err, "expected no error when deleting loadBalancer")

	loadBalancer, err := c.LoadBalancer.Get(projectID, loadBalancerID)
	require.Error(t, err, "expected error when getting loadBalancer after deletion")
	require.Nil(t, loadBalancer, "the loadBalancer should not exist anymore")
}
