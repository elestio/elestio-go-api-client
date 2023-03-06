package elestio

import (
	"fmt"
	"os"
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

	projectID := "3234"
	serviceID := "29648534"

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

func TestServiceHandler_Create(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectId := "2318"

	service, err := c.Service.Create(CreateServiceRequest{
		ProjectID:    projectId,
		ServerName:   "pg-asia",
		ServerType:   "MICRO-1C-1G",
		TemplateID:   11,
		Version:      "14",
		ProviderName: "Amazon Lightsail",
		Datacenter:   "ap-northeast-2",
		SupportLevel: "level1",
		AdminEmail:   "adamkrim.dev@gmail.com",
	})

	require.NoError(t, err, "expected no error when creating service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, "pg-asia", service.ServerName, "expected service name to be pg-asia")
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

	serviceID := "26454028"
	projectID := "2130"

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

func TestServiceHandler_EnableFirewall(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	serviceID := "26454028"
	projectID := "2130"

	service, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, service, "expected non-nil service")
	require.Equal(t, serviceID, service.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(0), service.FirewallEnabled, "expected firewall to be disabled")

	err = c.Service.EnableFirewall(serviceID)
	require.NoError(t, err, "expected no error when enabling firewall")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")
	require.NotNil(t, updatedService, "expected non-nil service")
	require.Equal(t, serviceID, updatedService.ID, "expected service ID to be"+serviceID)
	require.Equal(t, NumberAsBool(1), updatedService.FirewallEnabled, "expected firewall to be enabled")
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

func TestServiceHandler_AddSSHKey(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "3234"
	serviceID := "29648534"

	err := c.Service.AddSSHKey(serviceID, "test", "ssh-rsa fakeKey adam@macbook")
	require.NoError(t, err, "expected no error when adding ssh key")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")

	fmt.Fprintf(os.Stdout, "Service: %v", updatedService.SSHKeys)
}

func TestServiceHandler_RemoveSSHKey(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "3234"
	serviceID := "29648534"

	err := c.Service.RemoveSSHKey(serviceID, "test")
	require.NoError(t, err, "expected no error when removing ssh key")

	updatedService, err := c.Service.Get(projectID, serviceID)
	require.NoError(t, err, "expected no error when getting service")

	fmt.Fprintf(os.Stdout, "Service: %v", updatedService.SSHKeys)
}
