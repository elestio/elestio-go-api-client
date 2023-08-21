package elestio

import (
	"fmt"
	"strings"
)

// ServiceHandler is the client handler for service endpoints.
type ServiceHandler struct {
	client *Client
}

const (
	ServiceStatusStopped   string = "off"
	ServiceStatusRunning   string = "running"
	ServiceStatusDeleting  string = "deleting"
	ServiceStatusMigrating string = "migrating"

	ServiceDeploymentStatusDeployed   string = "Deployed"
	ServiceDeploymentStatusInProgress string = "IN PROGRESS"

	ServiceFirewallPortProtocolTCP ServiceFirewallPortProtocol = "tcp"
	ServiceFirewallPortProtocolUDP ServiceFirewallPortProtocol = "udp"
)

type (
	// NumberAsBool is a type to handle boolean values as numbers.
	// '0' is false. '1' is true.
	NumberAsBool uint8

	Template struct {
		ID                  int64  `json:"id"`
		Name                string `json:"title"`
		Category            string `json:"category"`
		Description         string `json:"description"`
		Logo                string `json:"mainImage"`
		DockerHubImage      string `json:"dockerhub_image"`
		DockerHubDefaultTag string `json:"dockerhub_default_tag"`
	}

	ServiceAdmin struct {
		URL      string `json:"url"`
		User     string `json:"user"`
		Password string `json:"password"`
	}

	ServiceDatabaseAdmin struct {
		Host     string
		Port     string
		User     string
		Password string
		Command  string
	}

	ServiceSSHPublicKey struct {
		ID   int64  `json:"indexID"`
		Name string `json:"name"`
		Key  string `json:"key"`
	}

	ServiceFirewallPort struct {
		// Port can be a single port or a range of ports. For example: "80" or "80-90".
		Port     string
		Protocol ServiceFirewallPortProtocol
	}

	ServiceFirewallPortProtocol string

	Service struct {
		ID                                          string                `json:"vmID"`
		ProjectID                                   string                `json:"projectID"`
		ServerName                                  string                `json:"displayName"`
		ServerType                                  string                `json:"serverType"`
		TemplateID                                  int64                 `json:"template"`
		Version                                     string                `json:"selected_software_tag"`
		ProviderName                                string                `json:"provider"`
		Datacenter                                  string                `json:"datacenter"`
		SupportLevel                                string                `json:"support"`
		AdminEmail                                  string                `json:"email"`
		Category                                    string                `json:"category"`
		Status                                      string                `json:"status"`
		DeploymentStatus                            string                `json:"deploymentStatus"`
		DeploymentStartedAt                         string                `json:"deploymentStart"`
		DeploymentEndedAt                           string                `json:"deploymentEnd"`
		CreatorName                                 string                `json:"creatorName"`
		CreatedAt                                   string                `json:"creationDate"`
		IPV4                                        string                `json:"ipv4"`
		IPV6                                        string                `json:"ipv6"`
		CNAME                                       string                `json:"cname"`
		CustomDomainNames                           []string              `json:"customDomainNames"`
		SSHPublicKeys                               []ServiceSSHPublicKey `json:"sshKeys"`
		Country                                     string                `json:"country"`
		City                                        string                `json:"city"`
		AdminUser                                   string                `json:"adminUser"`
		RootAppPath                                 string                `json:"rootAppPath"`
		GlobalIP                                    string                `json:"globalIP"`
		AdminExternalPort                           int64                 `json:"adminExternalPort"`
		AdminInternalPort                           int64                 `json:"adminInternalPort"`
		TrafficOutgoing                             int64                 `json:"traffic_outgoing"`
		TrafficIncoming                             int64                 `json:"traffic_ingoing"`
		TrafficIncluded                             int64                 `json:"traffic_included"`
		Cores                                       int64                 `json:"cores"`
		RAMSizeGB                                   string                `json:"ramGB"`
		StorageSizeGB                               int64                 `json:"storageSizeGB"`
		PricePerHour                                string                `json:"pricePerHour"`
		AppAutoUpdatesEnabled                       NumberAsBool          `json:"app_AutoUpdate_Enabled"`
		AppAutoUpdatesDayOfWeek                     int64                 `json:"app_AutoUpdate_DayOfWeek"`
		AppAutoUpdatesHour                          int64                 `json:"app_AutoUpdate_Hour"`
		AppAutoUpdatesMinute                        int64                 `json:"app_AutoUpdate_Minute"`
		SystemAutoUpdatesEnabled                    NumberAsBool          `json:"system_AutoUpdate_Enabled"`
		SystemAutoUpdatesSecurityPatchesOnlyEnabled NumberAsBool          `json:"system_AutoUpdate_SecurityPatchesOnly"`
		SystemAutoUpdatesRebootDayOfWeek            int64                 `json:"system_AutoUpdate_RebootDayOfWeek"`
		SystemAutoUpdatesRebootHour                 int64                 `json:"system_AutoUpdate_RebootHour"`
		SystemAutoUpdatesRebootMinute               int64                 `json:"system_AutoUpdate_RebootMinute"`
		BackupsEnabled                              NumberAsBool          `json:"backupsActivated"`
		RemoteBackupsEnabled                        NumberAsBool          `json:"remoteBackupsActivated"`
		ExternalBackupsEnabled                      NumberAsBool          `json:"isExternalBackupActivated"`
		ExternalBackupsUpdateDayOfWeek              int64                 `json:"externalBackupUpdateDay"`
		ExternalBackupsUpdateHour                   int64                 `json:"externalBackupUpdateHour"`
		ExternalBackupsUpdateMinute                 int64                 `json:"externalBackupUpdateMinute"`
		ExternalBackupsUpdateType                   string                `json:"externalBackupUpdateType"`
		ExternalBackupsRetainDayOfWeek              int64                 `json:"externalBackupRetainDay"`
		FirewallEnabled                             NumberAsBool          `json:"isFirewallActivated"`
		FirewallID                                  string                `json:"firewall_id"`
		FirewallPorts                               string                `json:"firewallPorts"`
		Env                                         map[string]string
		Admin                                       ServiceAdmin
		DatabaseAdmin                               ServiceDatabaseAdmin
		DatabaseAdminCommand                        string       `json:"managedDBCLI"`
		DatabaseAdminPort                           string       `json:"managedDBPort"`
		AlertsEnabled                               NumberAsBool `json:"isAlertsActivated"`
	}

	CreateServiceRequest struct {
		ProjectID                 string       `json:"projectId"`
		ServerName                string       `json:"serverName"`
		ServerType                string       `json:"serverType"`
		TemplateID                int64        `json:"templateID"`
		Version                   string       `json:"version"`
		ProviderName              string       `json:"providerName"`
		Datacenter                string       `json:"datacenter"`
		SupportLevel              string       `json:"support"`
		AppPassword               string       `json:"appPassword"`
		AdminEmail                string       `json:"adminEmail"`
		IsSystemAutoUpdateEnabled NumberAsBool `json:"system_AutoUpdate_Enabled"`
		IsAppAutoUpdateEnabled    NumberAsBool `json:"app_AutoUpdate_Enabled"`
	}
)

func (h *ServiceHandler) GetTemplatesList() ([]*Template, error) {
	type getTemplatesListResponse struct {
		Templates []Template `json:"instances"`
	}

	bts, err := h.client.sendGetRequest(fmt.Sprintf("%s/api/servers/getTemplates", h.client.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	var res getTemplatesListResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	if res.Templates == nil || len(res.Templates) == 0 {
		return nil, fmt.Errorf("templates not found")
	}

	var templates []*Template
	for _, template := range res.Templates {
		template := template // avoid iteration with same pointer
		template.Logo = strings.Replace(template.Logo, "//", "https://", 1)
		templates = append(templates, &template)
	}

	return templates, nil
}

func (h *ServiceHandler) Get(projectID, serviceID string) (*Service, error) {
	type getServiceRequest struct {
		ProjectID string `json:"projectID"`
		ServiceID string `json:"vmID"`
		JWT       string `json:"jwt"`
	}

	type getServiceResponse struct {
		APIResponse
		Services []Service `json:"serviceInfos"`
	}

	req := getServiceRequest{
		ProjectID: projectID,
		ServiceID: serviceID,
		JWT:       h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/servers/getServerDetails", h.client.BaseURL),
		req,
	)
	if err != nil {
		return nil, err
	}

	var res getServiceResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	if res.Services == nil || len(res.Services) == 0 {
		return nil, fmt.Errorf("service not found")
	}

	return h.formatServiceForClient(&res.Services[0])
}

func (h *ServiceHandler) GetList(projectID string) ([]*Service, error) {
	type getListServiceRequest struct {
		ProjectID       string `json:"projectId"`
		AppID           string `json:"appid"`
		IsActiveService bool   `json:"isActiveService"`
		JWT             string `json:"jwt"`
	}

	type getListServiceResponse struct {
		APIResponse
		Services []Service `json:"servers"`
	}

	req := getListServiceRequest{
		ProjectID:       projectID,
		AppID:           "",
		IsActiveService: true,
		JWT:             h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/servers/getServices", h.client.BaseURL),
		req,
	)
	if err != nil {
		return nil, err
	}

	var res getListServiceResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	var services []*Service
	for i := range res.Services {
		s, err := h.formatServiceForClient(&res.Services[i])
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}

func (h *ServiceHandler) Create(req CreateServiceRequest) (*Service, error) {
	type createServiceFullRequest struct {
		CreateServiceRequest
		Data                  string `json:"data"`
		AppID                 string `json:"appid"`
		DeploymentServiceType string `json:"deploymentServiceType"` // "normal"
		ServiceType           string `json:"serviceType"`           // "service"
		JWT                   string `json:"jwt"`
	}

	fullReq := createServiceFullRequest{
		CreateServiceRequest:  req,
		Data:                  "",
		AppID:                 "",
		DeploymentServiceType: "normal",
		ServiceType:           "Service",
		JWT:                   h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/servers/createServer", h.client.BaseURL),
		fullReq,
	)
	if err != nil {
		return nil, err
	}

	type createServiceResponse struct {
		APIResponse
		ID   []FlexString `json:"providerServerID"`
		Data []Service    `json:"data"`
	}

	var res createServiceResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	return h.formatServiceForClient(&res.Data[0])
}

func (h *ServiceHandler) Delete(projectID, serviceID string, keepBackups bool) error {
	type deleteServiceRequest struct {
		ProjectID       string `json:"projectID"`
		ServiceID       string `json:"vmID"`
		IsWithoutBackup bool   `json:"isDeleteServiceWithBackup"`
		JWT             string `json:"jwt"`
	}

	type deleteServiceResponse struct {
		APIResponse
	}

	req := deleteServiceRequest{
		ProjectID:       projectID,
		ServiceID:       serviceID,
		IsWithoutBackup: !keepBackups,
		JWT:             h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/deleteServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	var res deleteServiceResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return err
	}

	return nil
}

func (h *ServiceHandler) UpdateVersion(serviceId string, newVersion string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
		Version   string `json:"versionTag"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    "softwareChangeSelectedVersion",
		Version:   newVersion,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

// UpdateServerType updates the server type of a service.
// You can only upgrade the server type, not downgrade.
// The service will reboot in a few minutes.
func (h *ServiceHandler) UpdateServerType(serviceId string, newServerType string, providerName string, datacenter string) error {
	req := struct {
		JWT               string `json:"jwt"`
		ServiceID         string `json:"vmID"`
		Action            string `json:"action"`
		ServerType        string `json:"newType"`
		ProviderName      string `json:"providerName"`
		Datacenter        string `json:"region"`
		UpgradeCPURAMOnly bool   `json:"upgradeCPURAMOnly"`
	}{
		JWT:               h.client.jwt,
		ServiceID:         serviceId,
		Action:            "changeType",
		ServerType:        newServerType,
		ProviderName:      providerName,
		Datacenter:        datacenter,
		UpgradeCPURAMOnly: false,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) DisableAppAutoUpdates(serviceId string) error {
	return h.DoActionOnServer(serviceId, "appAutoUpdateDisable")
}

func (h *ServiceHandler) EnableAppAutoUpdates(serviceId string) error {
	req := struct {
		JWT             string `json:"jwt"`
		ServiceID       string `json:"vmID"`
		Action          string `json:"action"`
		UpdateDayOfWeek string `json:"appAutoUpdateDayOfWeek"`
		UpdateHour      string `json:"appAutoUpdateHour"`
		UpdateMinute    string `json:"appAutoUpdateMinute"`
	}{
		JWT:             h.client.jwt,
		ServiceID:       serviceId,
		Action:          "appAutoUpdateEnable",
		UpdateDayOfWeek: "0",
		UpdateHour:      "1",
		UpdateMinute:    "00",
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) DisableSystemAutoUpdates(serviceId string) error {
	return h.DoActionOnServer(serviceId, "systemAutoUpdateDisable")
}

func (h *ServiceHandler) EnableSystemAutoUpdates(serviceId string, isSystemAutoUpdatesSecurityPatchesOnlyEnabled bool) error {
	req := struct {
		JWT                                           string `json:"jwt"`
		ServiceID                                     string `json:"vmID"`
		Action                                        string `json:"action"`
		UpdateDayOfWeek                               string `json:"systemAutoUpdateRebootDayOfWeek"`
		UpdateHour                                    string `json:"systemAutoUpdateRebootHour"`
		UpdateMinute                                  string `json:"systemAutoUpdateRebootMinute"`
		IsSystemAutoUpdatesSecurityPatchesOnlyEnabled bool   `json:"systemAutoUpdateSecurityPatchesOnly"`
	}{
		JWT:             h.client.jwt,
		ServiceID:       serviceId,
		Action:          "systemAutoUpdateEnable",
		UpdateDayOfWeek: "0",
		UpdateHour:      "5",
		UpdateMinute:    "00",
		IsSystemAutoUpdatesSecurityPatchesOnlyEnabled: isSystemAutoUpdatesSecurityPatchesOnlyEnabled,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) DisableBackups(serviceId string) error {
	return h.DoActionOnServer(serviceId, "disableBackup")
}

func (h *ServiceHandler) EnableBackups(serviceId string) error {
	return h.DoActionOnServer(serviceId, "enableBackup")
}

func (h *ServiceHandler) DisableRemoteBackups(serviceId string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"serverID"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/backups/DisableAutoBackups", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) EnableRemoteBackups(serviceId string) error {
	req := struct {
		JWT        string `json:"jwt"`
		ServiceID  string `json:"serverID"`
		BackupPath string `json:"backupPath"`
		BackupHour int64  `json:"backupHour"`
	}{
		JWT:        h.client.jwt,
		ServiceID:  serviceId,
		BackupPath: "/opt",
		BackupHour: 4,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/backups/SetupAutoBackups", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) DisableAlerts(serviceId string) error {
	return h.DoActionOnServer(serviceId, "disableAlerts")
}

func (h *ServiceHandler) EnableAlerts(serviceId string) error {
	req := struct {
		JWT                 string `json:"jwt"`
		ServiceID           string `json:"vmID"`
		Action              string `json:"action"`
		MonitCycleInSeconds int64  `json:"monitCycleInSeconds"`
		Rules               string `json:"rules"`
	}{
		JWT:                 h.client.jwt,
		ServiceID:           serviceId,
		Action:              "enableAlerts",
		MonitCycleInSeconds: 60,
		Rules:               "[{\"parameter\":\"CPU\",\"value\":90,\"cycles\":15,\"unit\":\"%\"},{\"parameter\":\"MEMORY\",\"value\":90,\"cycles\":15,\"unit\":\"%\"},{\"parameter\":\"SWAP\",\"value\":75,\"cycles\":15,\"unit\":\"%\"},{\"parameter\":\"SPACE\",\"value\":80,\"cycles\":15,\"unit\":\"%\"},{\"parameter\":\"INODE\",\"value\":80,\"cycles\":15,\"unit\":\"%\"},{\"parameter\":\"READ_RATE\",\"value\":20,\"cycles\":15,\"unit\":\"MB/s\"},{\"parameter\":\"WRITE_RATE\",\"value\":20,\"cycles\":15,\"unit\":\"MB/s\"},{\"parameter\":\"SATURATION\",\"value\":90,\"cycles\":15,\"unit\":\"%\"},{\"parameter\":\"DOWNLOAD\",\"value\":25,\"cycles\":15,\"unit\":\"MB/s\"},{\"parameter\":\"UPLOAD\",\"value\":25,\"cycles\":15,\"unit\":\"MB/s\"}]",
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) DisableFirewall(serviceId string) error {
	return h.DoActionOnServer(serviceId, "disableFirewall")
}

func (h *ServiceHandler) EnableFirewall(serviceId string, ports []ServiceFirewallPort) error {
	acceptedProtocols := []ServiceFirewallPortProtocol{ServiceFirewallPortProtocolTCP, ServiceFirewallPortProtocolUDP}
	portMap := make(map[string]bool)
	for _, port := range ports {
		portMap[port.Port] = true

		if !Contains(acceptedProtocols, port.Protocol) {
			return fmt.Errorf("invalid firewall port protocol %s", port.Protocol)
		}
	}

	if _, exists := portMap["22"]; !exists {
		ports = append(ports, ServiceFirewallPort{
			Port:     "22",
			Protocol: ServiceFirewallPortProtocolTCP,
		})
	}

	if _, exists := portMap["4242"]; !exists {
		ports = append(ports, ServiceFirewallPort{
			Port:     "4242",
			Protocol: ServiceFirewallPortProtocolUDP,
		})
	}

	var rules []string

	for _, port := range ports {
		rules = append(rules, fmt.Sprintf("{\"type\":\"INPUT\",\"port\":\"%s\",\"protocol\":\"%s\",\"targets\":[\"0.0.0.0/0\",\"::/0\"]}", port.Port, port.Protocol))
	}

	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
		Rules     string `json:"rules"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    "enableFirewall",
		Rules:     fmt.Sprintf("[%s]", strings.Join(rules, ",")),
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) AddCustomDomainName(serviceId string, domain string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
		Domain    string `json:"domain"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    "SSLDomainsAdd",
		Domain:    domain,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) RemoveCustomDomainName(serviceId string, domain string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
		Domain    string `json:"domain"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    "SSLDomainsRemove",
		Domain:    domain,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) AddSSHPublicKey(serviceId string, name string, key string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
		Name      string `json:"name"`
		Key       string `json:"key"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    "SSHPubKeysAdd",
		Name:      name,
		Key:       key,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) RemoveSSHPublicKey(serviceId string, name string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
		Name      string `json:"deleteParams"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    "SSHPubKeysRemove",
		Name:      name,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) RebootServer(serviceId string) error {
	return h.DoActionOnServer(serviceId, "reboot")
}

func (h *ServiceHandler) DoActionOnServer(serviceId string, action string) error {
	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
	}{
		JWT:       h.client.jwt,
		ServiceID: serviceId,
		Action:    action,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	return checkAPIResponse(bts, nil)
}

func (h *ServiceHandler) GetServiceEnv(service *Service) (*map[string]string, error) {
	envMap, emptyEnvMap := make(map[string]string), make(map[string]string)

	if service.DeploymentStatus != ServiceDeploymentStatusDeployed {
		return &emptyEnvMap, nil
	}

	req := struct {
		JWT        string `json:"jwt"`
		ProjectID  string `json:"projectID"`
		ServiceID  string `json:"vmID"`
		TemplateID int64  `json:"templateID"`
		Action     string `json:"action"`
	}{
		JWT:        h.client.jwt,
		ProjectID:  service.ProjectID,
		ServiceID:  service.ID,
		TemplateID: service.TemplateID,
		Action:     "getAppStackConfig",
	}

	res := struct {
		APIResponse
		Data struct {
			Env string `json:"envResult"`
		} `json:"data"`
	}{}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return &emptyEnvMap, nil
	}

	if err := checkAPIResponse(bts, &res); err != nil {
		return &emptyEnvMap, nil
	}

	envs := strings.Split(res.Data.Env, "\n")
	for _, env := range envs {
		envSplit := strings.Split(env, "=")
		if len(envSplit) == 2 {
			envMap[envSplit[0]] = envSplit[1]
		}
	}

	return &envMap, nil
}

// GetServiceAdmin returns the admin credentials for a service,
// returns an empty ServiceAdmin if the service is not deployed.
func (h *ServiceHandler) GetServiceAdmin(service *Service) (*ServiceAdmin, error) {
	serviceAdmin, emptyServiceAdmin := ServiceAdmin{}, ServiceAdmin{}

	if service.DeploymentStatus != ServiceDeploymentStatusDeployed {
		return &emptyServiceAdmin, nil
	}

	req := struct {
		JWT               string `json:"jwt"`
		ProjectID         string `json:"projectID"`
		ServiceID         string `json:"vmID"`
		AppId             string `json:"appId"`
		IsServerDeleted   bool   `json:"isServerDeleted"`
		AdminExternalPort int64  `json:"srvPort"`
		AdminInternalPort int64  `json:"targetPort"`
	}{
		JWT:               h.client.jwt,
		ProjectID:         service.ProjectID,
		ServiceID:         service.ID,
		AppId:             "CloudVM",
		IsServerDeleted:   false,
		AdminExternalPort: service.AdminExternalPort,
		AdminInternalPort: service.AdminInternalPort,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/getAppCredentials", h.client.BaseURL), req)
	if err != nil {
		return &emptyServiceAdmin, nil
	}

	if err := checkAPIResponse(bts, &serviceAdmin); err != nil {
		return &emptyServiceAdmin, nil
	}

	return &serviceAdmin, nil
}

// GetServiceDatabaseAdmin returns the database admin credentials for a service,
// returns an empty ServiceDatabaseAdmin if the service is not deployed.
func (h *ServiceHandler) GetServiceDatabaseAdmin(service *Service) (*ServiceDatabaseAdmin, error) {
	databaseAdmin, emptyDatabaseAdmin := ServiceDatabaseAdmin{}, ServiceDatabaseAdmin{}

	if service.DeploymentStatus != ServiceDeploymentStatusDeployed {
		return &emptyDatabaseAdmin, nil
	}

	if service.Category != "Databases & Cache" {
		return &emptyDatabaseAdmin, nil
	}

	req := struct {
		JWT               string `json:"jwt"`
		ProjectID         string `json:"projectID"`
		ServiceID         string `json:"vmID"`
		AppId             string `json:"appId"`
		IsServerDeleted   bool   `json:"isServerDeleted"`
		AdminExternalPort int64  `json:"srvPort"`
		AdminInternalPort int64  `json:"targetPort"`
		Mode              string `json:"mode"`
	}{
		JWT:               h.client.jwt,
		ProjectID:         service.ProjectID,
		ServiceID:         service.ID,
		AppId:             "CloudVM",
		IsServerDeleted:   false,
		AdminExternalPort: service.AdminExternalPort,
		AdminInternalPort: service.AdminInternalPort,
		Mode:              "dbAdmin",
	}

	res := struct{ ServiceAdmin }{}
	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/getAppCredentials", h.client.BaseURL), req)
	if err != nil {
		return &emptyDatabaseAdmin, nil
	}

	if err := checkAPIResponse(bts, &res); err != nil {
		return &emptyDatabaseAdmin, nil
	}

	databaseAdmin.Host = service.CNAME
	databaseAdmin.Port = service.DatabaseAdminPort
	databaseAdmin.User = res.User
	databaseAdmin.Password = res.Password
	databaseAdmin.Command = service.DatabaseAdminCommand
	databaseAdmin.Command = strings.Replace(databaseAdmin.Command, "[APP_PASSWORD]", databaseAdmin.Password, -1)
	databaseAdmin.Command = strings.Replace(databaseAdmin.Command, "[EMAIL]", service.AdminEmail, -1)
	databaseAdmin.Command = strings.Replace(databaseAdmin.Command, "[DOMAIN]", databaseAdmin.Host, -1)

	return &databaseAdmin, nil
}

// GetServiceCustomDomainNames returns the custom domain names configured for a service
func (h *ServiceHandler) GetServiceCustomDomainNames(service *Service) (*[]string, error) {
	var empty []string

	if service.DeploymentStatus != ServiceDeploymentStatusDeployed {
		return &empty, nil
	}

	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
	}{
		JWT:       h.client.jwt,
		ServiceID: service.ID,
		Action:    "SSLDomainsList",
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return &empty, nil
	}

	var customDomainNames []string
	if err := checkAPIResponse(bts, &customDomainNames); err != nil {
		return &empty, nil
	}

	// Remove the default service CNAME from the list of custom domain names
	customDomainNames = RemoveStringFromSlice(customDomainNames, service.CNAME)

	return &customDomainNames, nil
}

// GetServiceSSHPublicKeys returns the ssh public keys configured for a service
func (h *ServiceHandler) GetServiceSSHPublicKeys(service *Service) (*[]ServiceSSHPublicKey, error) {
	var empty []ServiceSSHPublicKey

	if service.DeploymentStatus != ServiceDeploymentStatusDeployed {
		return &empty, nil
	}

	req := struct {
		JWT       string `json:"jwt"`
		ServiceID string `json:"vmID"`
		Action    string `json:"action"`
	}{
		JWT:       h.client.jwt,
		ServiceID: service.ID,
		Action:    "SSHPubKeysList",
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), req)
	if err != nil {
		return &empty, nil
	}

	res := struct {
		APIResponse
		Data []ServiceSSHPublicKey `json:"data"`
	}{}

	if err := checkAPIResponse(bts, &res); err != nil {
		return &empty, nil
	}

	return &res.Data, nil
}

func (h *ServiceHandler) formatServiceForClient(service *Service) (*Service, error) {
	if service == nil {
		return nil, fmt.Errorf("cannot format nil service")
	}

	service.AdminUser = strings.Replace(service.AdminUser, "[EMAIL]", service.AdminEmail, -1)

	env, err := h.GetServiceEnv(service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service env: %s", err)
	}
	service.Env = *env

	admin, err := h.GetServiceAdmin(service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service admin: %s", err)
	}
	service.Admin = *admin

	databaseAdmin, err := h.GetServiceDatabaseAdmin(service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service database admin: %s", err)
	}
	service.DatabaseAdmin = *databaseAdmin

	customDomainNames, err := h.GetServiceCustomDomainNames(service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service custom domain names: %s", err)
	}
	service.CustomDomainNames = *customDomainNames

	sshPublicKeys, err := h.GetServiceSSHPublicKeys(service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service ssh public keys: %s", err)
	}
	service.SSHPublicKeys = *sshPublicKeys

	return service, nil
}
