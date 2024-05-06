package elestio

import (
	"fmt"
)

type LoadBalancerHandler struct {
	client *Client
}

const (
	LoadBalancerDeploymentStatusDeployed   string = "Deployed"
	LoadBalancerDeploymentStatusInProgress string = "IN PROGRESS"
)

type (
	LoadBalancer struct {
		ID               string
		ProjectID        string
		ProviderName     string
		Datacenter       string
		ServerType       string
		Config           LoadBalancerConfig
		CreatedAt        string
		CreatorName      string
		DeploymentStatus string
		IPV4             string
		IPV6             string
		CNAME            string
		Country          string
		City             string
		GlobalIP         string
		Cores            int64
		RAMSizeGB        string
		StorageSizeGB    int64
		PricePerHour     string
	}

	LoadBalancerConfig struct {
		HostHeader             string
		IsAccessLogsEnabled    bool
		IsForceHTTPSEnabled    bool
		IPRateLimit            int64
		IsIPRateLimitEnabled   bool
		OutputCacheInSeconds   int64
		IsStickySessionEnabled bool
		IsProxyProtocolEnabled bool
		SSLDomains             []string
		ForwardRules           []LoadBalancerConfigForwardRule
		OutputHeaders          []LoadBalancerConfigOutputHeader
		TargetServices         []string
		RemoveResponseHeaders  []string
	}

	LoadBalancerConfigForwardRule struct {
		Protocol       string `json:"protocol"`
		Port           string `json:"listeningPort"`
		TargetProtocol string `json:"targetProtocol"`
		TargetPort     string `json:"targetPort"`
	}

	LoadBalancerConfigOutputHeader struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

func (h *LoadBalancerHandler) Get(projectID, loadBalancerID string) (*LoadBalancer, error) {
	// Fetch load balancer details
	reqDetails := struct {
		ProjectID      string `json:"projectID"`
		LoadBalancerID string `json:"vmID"`
		JWT            string `json:"jwt"`
	}{
		ProjectID:      projectID,
		LoadBalancerID: loadBalancerID,
		JWT:            h.client.jwt,
	}
	btsDetails, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/servers/getServerDetails", h.client.BaseURL),
		reqDetails,
	)
	if err != nil {
		return nil, err
	}
	var resDetails struct {
		APIResponse
		Services []struct {
			CreatedAt        string `json:"creationDate"`
			CreatorName      string `json:"creatorName"`
			DeploymentStatus string `json:"deploymentStatus"`
			IPV4             string `json:"ipv4"`
			IPV6             string `json:"ipv6"`
			CNAME            string `json:"cname"`
			Country          string `json:"country"`
			City             string `json:"city"`
			GlobalIP         string `json:"globalIP"`
			Cores            int64  `json:"cores"`
			RAMSizeGB        string `json:"ramGB"`
			StorageSizeGB    int64  `json:"storageSizeGB"`
			PricePerHour     string `json:"pricePerHour"`
			// TODO: Add all API available fields later
			// SSHKeys          []ServiceSSHKey `json:"sshKeys"`
			// AlertsEnabled    NumberAsBool    `json:"isAlertsActivated"`
		} `json:"serviceInfos"`
	}
	if err = checkAPIResponse(btsDetails, &resDetails); err != nil {
		return nil, err
	}
	// API returns an array of services, but we only need the first one
	if len(resDetails.Services) == 0 {
		return nil, fmt.Errorf("load balancer not found")
	}
	details := resDetails.Services[0]

	// Fetch load balancer config
	reqConfig := struct {
		ProjectID       string `json:"projectID"`
		LoadBalancerID  string `json:"loadBalancerID"`
		IsRestoreLb     bool   `json:"isRestoreLb"`
		IsActiveService bool   `json:"isActiveService"`
		JWT             string `json:"jwt"`
	}{
		ProjectID:       projectID,
		LoadBalancerID:  loadBalancerID,
		IsRestoreLb:     false,
		IsActiveService: true,
		JWT:             h.client.jwt,
	}
	btsConfig, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/loadBalancer/getLBDetails", h.client.BaseURL),
		reqConfig,
	)
	if err != nil {
		return nil, err
	}
	var resConfig struct {
		APIResponse
		Data struct {
			ProjectID              string                           `json:"projectID"`
			ProviderName           string                           `json:"providerName"`
			Datacenter             string                           `json:"providerRegion"`
			ServerType             string                           `json:"planType"`
			HostHeader             string                           `json:"hostHeader"`
			IsAccessLogsEnabled    bool                             `json:"accessLog"`
			IsForceHTTPSEnabled    bool                             `json:"forceHttps"`
			IPRateLimit            int64                            `json:"ipRateLimit"`
			IsIPRateLimitEnabled   bool                             `json:"isIpRateLimiter"`
			OutputCacheInSeconds   int64                            `json:"outputCache"`
			IsStickySessionEnabled bool                             `json:"isStickySessions"`
			IsProxyProtocolEnabled bool                             `json:"proxyProtocol"`
			SSLDomains             []string                         `json:"sslDomains"`
			ForwardRules           []LoadBalancerConfigForwardRule  `json:"forwardingRules"`
			OutputHeaders          []LoadBalancerConfigOutputHeader `json:"outputHeaders"`
			TargetServices         []string                         `json:"targetServiceIDs"`
			RemoveResponseHeaders  []string                         `json:"removeResponseHeaders"`
		} `json:"data"`
	}
	if err = checkAPIResponse(btsConfig, &resConfig); err != nil {
		return nil, err
	}
	config := resConfig.Data

	// Build load balancer struct
	loadBalancer := LoadBalancer{
		ID:           loadBalancerID,
		ProjectID:    config.ProjectID,
		ProviderName: config.ProviderName,
		Datacenter:   config.Datacenter,
		ServerType:   config.ServerType,
		Config: LoadBalancerConfig{
			HostHeader:             config.HostHeader,
			IsAccessLogsEnabled:    config.IsAccessLogsEnabled,
			IsForceHTTPSEnabled:    config.IsForceHTTPSEnabled,
			IPRateLimit:            config.IPRateLimit,
			IsIPRateLimitEnabled:   config.IsIPRateLimitEnabled,
			OutputCacheInSeconds:   config.OutputCacheInSeconds,
			IsStickySessionEnabled: config.IsStickySessionEnabled,
			IsProxyProtocolEnabled: config.IsProxyProtocolEnabled,
			SSLDomains:             config.SSLDomains,
			ForwardRules:           config.ForwardRules,
			OutputHeaders:          config.OutputHeaders,
			TargetServices:         config.TargetServices,
			RemoveResponseHeaders:  config.RemoveResponseHeaders,
		},
		CreatedAt:        details.CreatedAt,
		CreatorName:      details.CreatorName,
		DeploymentStatus: details.DeploymentStatus,
		IPV4:             details.IPV4,
		IPV6:             details.IPV6,
		CNAME:            details.CNAME,
		Country:          details.Country,
		City:             details.City,
		GlobalIP:         details.GlobalIP,
		Cores:            details.Cores,
		RAMSizeGB:        details.RAMSizeGB,
		StorageSizeGB:    details.StorageSizeGB,
		PricePerHour:     details.PricePerHour,
	}

	return &loadBalancer, nil
}

type CreateLoadBalancerRequest struct {
	ProjectID    string                          `json:"projectId"`
	ProviderName string                          `json:"providerName"`
	Datacenter   string                          `json:"datacenter"`
	ServerType   string                          `json:"serverType"`
	Config       CreateLoadBalancerRequestConfig `json:"loadBalancerPayload"`
	CreatedFrom  string                          `json:"createdFrom"`
}

type CreateLoadBalancerRequestConfig struct {
	HostHeader             string                           `json:"hostHeader"`
	IsAccessLogsEnabled    bool                             `json:"accessLog"`
	IsForceHTTPSEnabled    bool                             `json:"forceHttps"`
	IPRateLimit            int64                            `json:"ipRateLimit"`
	IsIPRateLimitEnabled   bool                             `json:"isIpRateLimiter"`
	OutputCacheInSeconds   int64                            `json:"outputCache"`
	IsStickySessionEnabled bool                             `json:"stickySession"`
	IsProxyProtocolEnabled bool                             `json:"proxyProtocol"`
	SSLDomains             []string                         `json:"sslDomains"`
	ForwardRules           []LoadBalancerConfigForwardRule  `json:"forwardRules"`
	OutputHeaders          []LoadBalancerConfigOutputHeader `json:"outputHeaders"`
	TargetServices         []string                         `json:"targetServiceIDs"`
	RemoveResponseHeaders  []string                         `json:"removeResponseHeaders"`
}

func (h *LoadBalancerHandler) Create(req CreateLoadBalancerRequest) (*LoadBalancer, error) {
	if req.CreatedFrom == "" {
		req.CreatedFrom = "goClient"
	}

	fullReq := struct {
		CreateLoadBalancerRequest
		ServiceType string `json:"serviceType"`
		JWT         string `json:"jwt"`
	}{
		CreateLoadBalancerRequest: req,
		ServiceType:               "LB",
		JWT:                       h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/servers/createServer", h.client.BaseURL),
		fullReq,
	)
	if err != nil {
		return nil, err
	}

	var res struct {
		APIResponse
		ID []FlexString `json:"providerServerID"`
	}
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	return h.Get(req.ProjectID, (string)(res.ID[0]))
}

type UpdateLoadBalancerConfigRequest struct {
	HostHeader             string                           `json:"hostHeader"`
	IsAccessLogsEnabled    bool                             `json:"accessLog"`
	IsForceHTTPSEnabled    bool                             `json:"forceHttps"`
	IPRateLimit            int64                            `json:"ipRateLimit"`
	IsIPRateLimitEnabled   bool                             `json:"isIpRateLimiter"`
	OutputCacheInSeconds   int64                            `json:"outputCache"`
	IsStickySessionEnabled bool                             `json:"stickySession"`
	IsProxyProtocolEnabled bool                             `json:"proxyProtocol"`
	SSLDomains             []string                         `json:"sslDomains"`
	ForwardRules           []LoadBalancerConfigForwardRule  `json:"forwardRules"`
	OutputHeaders          []LoadBalancerConfigOutputHeader `json:"outputHeaders"`
	TargetServices         []string                         `json:"targetServiceIDs"`
	RemoveResponseHeaders  []string                         `json:"removeResponseHeaders"`
}

func (h *LoadBalancerHandler) UpdateConfig(projectID string, loadBalancerID string, req UpdateLoadBalancerConfigRequest) (*LoadBalancer, error) {
	fullReq := struct {
		UpdateLoadBalancerConfigRequest
		LoadBalancerID string `json:"vmID"`
		Action         string `json:"action"`
		JWT            string `json:"jwt"`
	}{
		UpdateLoadBalancerConfigRequest: req,
		LoadBalancerID:                  loadBalancerID,
		Action:                          "updateLBSetting",
		JWT:                             h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/DoActionOnServer", h.client.BaseURL), fullReq)
	if err != nil {
		return nil, err
	}

	var res struct {
		APIResponse
	}
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	return h.Get(projectID, loadBalancerID)
}

func (h *LoadBalancerHandler) Delete(projectID, loadBalancerID string, keepBackups bool) error {
	type deleteLoadBalancerRequest struct {
		ProjectID       string `json:"projectID"`
		LoadBalancerID  string `json:"vmID"`
		IsWithoutBackup bool   `json:"isDeleteServiceWithBackup"`
		JWT             string `json:"jwt"`
	}

	type deleteLoadBalancerResponse struct {
		APIResponse
	}

	req := deleteLoadBalancerRequest{
		ProjectID:       projectID,
		LoadBalancerID:  loadBalancerID,
		IsWithoutBackup: !keepBackups,
		JWT:             h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(fmt.Sprintf("%s/api/servers/deleteServer", h.client.BaseURL), req)
	if err != nil {
		return err
	}

	var res deleteLoadBalancerResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return err
	}

	return nil
}
