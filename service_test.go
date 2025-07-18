package elestio

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupServiceTestCase(t *testing.T) *Client {
	t.Log("Setup service test case")

	email, apiKey := os.Getenv("ELESTIO_INTEGRATION_EMAIL"), os.Getenv("ELESTIO_INTEGRATION_API_KEY")

	c, err := NewClient(email, apiKey)
	if err != nil {
		t.Fatalf("client authentication error: %s", err)
	}

	return c
}

func TestServiceHandler_GetTemplatesList(t *testing.T) {
	t.Skip("Skipping test")
	c := NewUnsignedClient()

	templates, err := c.Service.GetTemplatesList()
	require.NoError(t, err, "expected no error when getting templates")
	require.NotNil(t, templates, "expected non-nil templates")
	require.NotEqual(t, 0, len(templates), "expected at least 1 template")
	// for _, template := range templates {
	// 	template := *template
	// 	fmt.Fprintf(os.Stdout, "Template: %v", template)
	// }
}

func TestServiceHandler_Get(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "13194"
	serviceID := "42438730"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)

	// fmt.Fprintf(os.Stdout, "Service: %v", service)
}

func TestServiceHandler_GetList(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "2234"

	services, err := c.Service.GetList(projectID)
	require.NoError(t, err, "expected no error when getting services")
	require.NotNil(t, services, "expected non-nil services")
	// require.Equal(t, 1, len(services.Service), "expected 1 service")
	// require.Equal(t, 1, services.Service[0].ID, "expected service ID to be 1")
}

func TestServiceHandler_ValidateConfig(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	validConfigReq := ValidateConfigRequest{
		TemplateId:   11,
		ProviderName: "hetzner",
		Datacenter:   "fsn1",
		ServerType:   "SMALL-1C-2G",
	}

	isValid, err := c.Service.ValidateConfig(validConfigReq)
	require.NoError(t, err, "expected no error when validating config")
	require.True(t, isValid, "expected config to be valid")

	isTemplateValid, err := c.Service.ValidateConfig(ValidateConfigRequest{
		TemplateId:   999999,
		ProviderName: validConfigReq.ProviderName,
		Datacenter:   validConfigReq.Datacenter,
		ServerType:   validConfigReq.ServerType,
	})
	require.Error(t, err, "expected error when validating invalid template config")
	require.False(t, isTemplateValid, "expected config to be invalid")

	isProviderNameValid, err := c.Service.ValidateConfig(ValidateConfigRequest{
		TemplateId:   validConfigReq.TemplateId,
		ProviderName: "invalid-provider",
		Datacenter:   validConfigReq.Datacenter,
		ServerType:   validConfigReq.ServerType,
	})
	require.Error(t, err, "expected error when validating invalid provider name config")
	require.False(t, isProviderNameValid, "expected config to be invalid")

	isDatacenterValid, err := c.Service.ValidateConfig(ValidateConfigRequest{
		TemplateId:   validConfigReq.TemplateId,
		ProviderName: validConfigReq.ProviderName,
		Datacenter:   "invalid-datacenter",
		ServerType:   validConfigReq.ServerType,
	})
	require.Error(t, err, "expected error when validating invalid datacenter config")
	require.False(t, isDatacenterValid, "expected config to be invalid")

	isServerTypeValid, err := c.Service.ValidateConfig(ValidateConfigRequest{
		TemplateId:   validConfigReq.TemplateId,
		ProviderName: validConfigReq.ProviderName,
		Datacenter:   validConfigReq.Datacenter,
		ServerType:   "invalid-servertype",
	})
	require.Error(t, err, "expected error when validating invalid server type config")
	require.False(t, isServerTypeValid, "expected config to be invalid")
}

func TestServiceHandler_Create(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectId := "596"

	service, err := c.Service.Create(CreateServiceRequest{
		ProjectID:    projectId,
		ServerName:   "mypostgres",
		ServerType:   "SMALL-2C-2G",
		TemplateID:   11,
		Version:      "14",
		ProviderName: "scaleway",
		Datacenter:   "fr-par-1",
		SupportLevel: "level1",
		AppPassword:  "L0ngPassw0rd",
		// AdminEmail:   "adamkrim.dev@gmail.com",
		CreatedFrom: "terraform",
	})

	require.NoError(t, err, "expected no error when creating service")
	require.NotNil(t, service, "expected non-nil service")
	fmt.Fprintf(os.Stdout, "Service: %v", service)
	require.Equal(t, "mypostgres", service.ServerName, "expected service name to be mypostgres")
}

func TestServiceHandler_Delete(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "27363147"

	err := c.Service.Delete(projectID, serviceID, false)
	require.NoError(t, err, "expected no error when deleting service")

	service, err := c.Service.Get(projectID, serviceID)
	require.Error(t, err, "expected error when getting service")
	require.Nil(t, service, "expected nil service")
}

func TestServiceHandler_UpdateVersion(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	currentVersion := "14"
	newVersion := "13"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)
	require.Equal(t, currentVersion, service.Version, "expected version to be enabled")

	err = c.Service.UpdateVersion(serviceID, newVersion)
	require.NoError(t, err, "expected no error when updating version")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, newVersion, updatedService.Version, "expected version changed")
}

func TestServiceHandler_UpdateServerType(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26728910"
	projectID := "2088"
	currentServerType := "SMALL-1C-2G"
	newServerType := "MEDIUM-2C-4G"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.Equal(t, currentServerType, service.ServerType, "expected matching current server type")

	err = c.Service.UpdateServerType(serviceID, newServerType, service.ProviderName, service.Datacenter)
	require.NoError(t, err, "expected no error when updating server type")
}

func TestServiceHandler_DisableAppAutoUpdates(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)
	require.Equal(t, NumberAsBool(1), service.AppAutoUpdatesEnabled, "expected app auto updates to be enabled")

	err = c.Service.DisableAppAutoUpdates(serviceID)
	require.NoError(t, err, "expected no error when disabling app auto updates")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), updatedService.AppAutoUpdatesEnabled, "expected app auto updates to be disabled")
}

func TestServiceHandler_EnableAppAutoUpdates(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), service.AppAutoUpdatesEnabled, "expected app auto updates to be disabled")

	err = c.Service.EnableAppAutoUpdates(serviceID)
	require.NoError(t, err, "expected no error when enabling app auto updates")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(1), updatedService.AppAutoUpdatesEnabled, "expected app auto updates to be enabled")
}

func TestServiceHandler_DisableSystemAutoUpdates(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)
	require.Equal(t, NumberAsBool(1), service.SystemAutoUpdatesEnabled, "expected system auto updates to be enabled")

	err = c.Service.DisableSystemAutoUpdates(serviceID)
	require.NoError(t, err, "expected no error when disabling system auto updates")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), updatedService.SystemAutoUpdatesEnabled, "expected system auto updates to be disabled")
}

func TestServiceHandler_EnableSystemAutoUpdates(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), service.SystemAutoUpdatesEnabled, "expected system auto updates to be disabled")

	err = c.Service.EnableSystemAutoUpdates(serviceID, true)
	require.NoError(t, err, "expected no error when enabling system auto updates")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(1), updatedService.SystemAutoUpdatesEnabled, "expected system auto updates to be enabled")
}

func TestServiceHandler_DisableBackups(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)
	require.Equal(t, NumberAsBool(1), service.BackupsEnabled, "expected backups to be enabled")

	err = c.Service.DisableBackups(serviceID)
	require.NoError(t, err, "expected no error when disabling backups")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), updatedService.BackupsEnabled, "expected backups to be disabled")
}

func TestServiceHandler_EnableBackups(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), service.BackupsEnabled, "expected backups to be disabled")

	err = c.Service.EnableBackups(serviceID)
	require.NoError(t, err, "expected no error when enabling backups")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(1), updatedService.BackupsEnabled, "expected backups to be enabled")
}

func TestServiceHandler_DisableRemoteBackups(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)
	require.Equal(t, NumberAsBool(1), service.RemoteBackupsEnabled, "expected remote backups to be enabled")

	err = c.Service.DisableRemoteBackups(serviceID)
	require.NoError(t, err, "expected no error when disabling remote backups")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), updatedService.RemoteBackupsEnabled, "expected remote backups to be disabled")
}

func TestServiceHandler_EnableRemoteBackups(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), service.RemoteBackupsEnabled, "expected remote backups to be disabled")

	err = c.Service.EnableRemoteBackups(serviceID)
	require.NoError(t, err, "expected no error when enabling remote backups")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(1), updatedService.RemoteBackupsEnabled, "expected remote backups to be enabled")
}

func TestServiceHandler_DisableAlerts(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26361706"

	service, err := c.Service.Get("2088", serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be 26361706")
	require.Equal(t, NumberAsBool(1), service.AlertsEnabled, "expected alerts to be enabled")

	err = c.Service.DisableAlerts(serviceID)
	require.NoError(t, err, "expected no error when disabling alerts")

	updatedService, err := c.Service.Get("2088", serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be 26361706")
	require.Equal(t, NumberAsBool(0), updatedService.AlertsEnabled, "expected alerts to be disabled")
}

func TestServiceHandler_EnableAlerts(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26361706"

	service, err := c.Service.Get("2088", serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be 26361706")
	require.Equal(t, NumberAsBool(0), service.AlertsEnabled, "expected alerts to be disabled")

	err = c.Service.EnableAlerts(serviceID)
	require.NoError(t, err, "expected no error when enabling alerts")

	updatedService, err := c.Service.Get("2088", serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be 26361706")
	require.Equal(t, NumberAsBool(1), updatedService.AlertsEnabled, "expected alerts to be enabled")
}

func TestServiceHandler_DisableFirewall(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "66470070"
	projectID := "596"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be "+serviceID)
	require.Equal(t, NumberAsBool(1), service.FirewallEnabled, "expected firewall to be enabled")

	err = c.Service.DisableFirewall(serviceID)
	require.NoError(t, err, "expected no error when disabling firewall")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), updatedService.FirewallEnabled, "expected firewall to be disabled")
}

func TestServiceHandler_EnableFirewallWithRules(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "66470070"
	projectID := "596"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), service.FirewallEnabled, "expected firewall to be disabled")

	// Test with valid rules including required ports
	rules := []ServiceFirewallRule{
		{
			Type:     ServiceFirewallRuleTypeInput,
			Port:     "22",
			Protocol: ServiceFirewallRuleProtocolTCP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
		{
			Type:     ServiceFirewallRuleTypeInput,
			Port:     "4242",
			Protocol: ServiceFirewallRuleProtocolUDP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
		{
			Type:     ServiceFirewallRuleTypeInput,
			Port:     "80",
			Protocol: ServiceFirewallRuleProtocolTCP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
		{
			Type:     ServiceFirewallRuleTypeInput,
			Port:     "443",
			Protocol: ServiceFirewallRuleProtocolTCP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
		{
			Type:     ServiceFirewallRuleTypeOutput,
			Port:     "443",
			Protocol: ServiceFirewallRuleProtocolTCP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
	}

	err = c.Service.EnableFirewallWithRules(serviceID, rules)
	require.NoError(t, err, "expected no error when enabling firewall with rules")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(1), updatedService.FirewallEnabled, "expected firewall to be enabled")
}

func TestServiceHandler_EnableFirewallWithRules_ValidationErrors(t *testing.T) {
	// t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "66470070"

	// Test invalid rule type
	invalidRulesType := []ServiceFirewallRule{
		{
			Type:     ServiceFirewallRuleTypeInput,
			Port:     "80",
			Protocol: ServiceFirewallRuleProtocolTCP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
		{
			Type:     "INVALID",
			Port:     "443",
			Protocol: ServiceFirewallRuleProtocolTCP,
			Targets:  []string{"0.0.0.0/0", "::/0"},
		},
	}

	err := c.Service.EnableFirewallWithRules(serviceID, invalidRulesType)
	require.Error(t, err, "expected error when using invalid rule type")
	require.Contains(t, err.Error(), "invalid rule type 'INVALID'", "expected specific error message about invalid rule type")
}

func TestServiceHandler_GetServiceFirewallRules(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "66470070"

	// Get the service first
	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")

	// Test getting firewall rules for a deployed service
	if service.DeploymentStatus == ServiceDeploymentStatusDeployed {
		rules, err := c.Service.GetServiceFirewallRules(service)
		require.NoError(t, err, "expected no error when getting firewall rules")
		require.NotNil(t, rules, "expected non-nil firewall rules")

		// Verify that rules are properly formatted
		for _, rule := range *rules {
			require.True(t, rule.Type == ServiceFirewallRuleTypeInput || rule.Type == ServiceFirewallRuleTypeOutput, "expected rule type to be INPUT or OUTPUT")
			require.NotEmpty(t, rule.Port, "expected rule port to not be empty")
			require.NotEmpty(t, rule.Protocol, "expected rule protocol to not be empty")
			require.NotNil(t, rule.Targets, "expected rule targets to not be nil")
		}

		// If firewall is enabled, we should have some rules
		if service.FirewallEnabled == NumberAsBool(1) {
			require.NotEmpty(t, *rules, "expected at least some firewall rules when firewall is enabled")
		}

		fmt.Fprintf(os.Stdout, "Firewall rules: %v", rules)
	}
}

func TestServiceHandler_GetServiceFirewallRules_NotDeployed(t *testing.T) {
	nonDeployedService := &Service{
		ID:               "test-service",
		ProjectID:        "test-project",
		DeploymentStatus: ServiceDeploymentStatusInProgress,
	}

	c := NewUnsignedClient()
	rules, err := c.Service.GetServiceFirewallRules(nonDeployedService)
	require.NoError(t, err, "expected no error when getting firewall rules for non-deployed service")
	require.NotNil(t, rules, "expected non-nil firewall rules")
	require.Empty(t, *rules, "expected empty firewall rules for non-deployed service")
}

func TestServiceHandler_GetServiceFirewallRules_FirewallDisabled(t *testing.T) {
	firewallDisabledService := &Service{
		ID:               "test-service",
		ProjectID:        "test-project",
		DeploymentStatus: ServiceDeploymentStatusDeployed,
		FirewallEnabled:  NumberAsBool(0), // Firewall disabled
	}

	c := NewUnsignedClient()
	rules, err := c.Service.GetServiceFirewallRules(firewallDisabledService)
	require.NoError(t, err, "expected no error when getting firewall rules for service with disabled firewall")
	require.NotNil(t, rules, "expected non-nil firewall rules")
	require.Empty(t, *rules, "expected empty firewall rules when firewall is disabled")
}

func TestServiceHandler_AddCustomDomain(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "28926765"

	err := c.Service.AddCustomDomainName(serviceID, "test.com")
	require.NoError(t, err, "expected no error when adding custom domain name")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")

	fmt.Fprintf(os.Stdout, "Service: %v", updatedService)
}

func TestServiceHandler_RemoveCustomDomain(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "28926765"

	err := c.Service.RemoveCustomDomainName(serviceID, "test.com")
	require.NoError(t, err, "expected no error when removing custom domain name")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")

	fmt.Fprintf(os.Stdout, "Service: %v", updatedService)
}

func TestServiceHandler_GetCustomDomainNames(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "103550055"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")

	customDomainNames, err := c.Service.GetServiceCustomDomainNames(service)
	require.NoError(t, err, "expected no error when getting custom domain names")
	require.NotNil(t, customDomainNames, "expected non-nil custom domain names")

	fmt.Fprintf(os.Stdout, "Custom domain names: %v", *customDomainNames)
}

func TestServiceHandler_AddSSHPublicKey(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "c4686e74-c75c-4ca8-9aaa-26f83eaaae97"
	keyName := "test"
	keyData := "ssh-rsa fakeKey test@macbook"

	err := c.Service.AddSSHPublicKey(serviceID, keyName, keyData)
	require.NoError(t, err, "expected no error when adding ssh key")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotEqual(
		t,
		-1,
		slices.IndexFunc(updatedService.SSHPublicKeys, func(s ServiceSSHPublicKey) bool { return s.Name == keyName && s.Key == keyData }),
		"expected ssh key to be added",
	)

	fmt.Fprintf(os.Stdout, "Service ssh public keys after added: %v", updatedService.SSHPublicKeys)
}

func TestServiceHandler_RemoveSSHPublicKey(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "596"
	serviceID := "c4686e74-c75c-4ca8-9aaa-26f83eaaae97"
	keyName := "test"
	keyData := "ssh-rsa fakeKey test@macbook"

	err := c.Service.RemoveSSHPublicKey(serviceID, "test")
	require.NoError(t, err, "expected no error when removing ssh key")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.Equal(
		t,
		-1,
		slices.IndexFunc(updatedService.SSHPublicKeys, func(s ServiceSSHPublicKey) bool { return s.Name == keyName && s.Key == keyData }),
		"expected ssh key to be removed",
	)

	fmt.Fprintf(os.Stdout, "Service ssh public key after remove: %v", updatedService.SSHPublicKeys)
}

func TestServiceHandler_Reboot(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "6450"
	serviceID := "409c54bf-e27e-4c5e-b760-aa0a3fcb2336"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting initial service")
	require.Equal(t, ServiceStatusRunning, service.Status, "expected initial service to be running")

	err = c.Service.RebootServer(serviceID)
	require.NoError(t, err, "expected no error when rebooting service")

	rebootedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting rebooted service")
	require.NotEqual(t, ServiceStatusRunning, rebootedService.Status, "expected rebooted service to be not running")
}
