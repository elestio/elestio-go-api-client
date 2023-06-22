package elestio

import (
	"fmt"
)

type LoadBalancerHandler struct {
	client *Client
}

type (
	LoadBalancer struct {
		ID           string
		ProjectID    string
		ProviderName string
		Datacenter   string
		ServerType   string
		Config       LoadBalancerConfig
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
	reqDetails := struct {
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

	btsDetails, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/loadBalancer/getLBDetails", h.client.BaseURL),
		reqDetails,
	)
	if err != nil {
		return nil, err
	}

	var resDetails struct {
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
	if err = checkAPIResponse(btsDetails, &resDetails); err != nil {
		return nil, err
	}

	loadBalancer := LoadBalancer{
		ID:           loadBalancerID,
		ProjectID:    resDetails.Data.ProjectID,
		ProviderName: resDetails.Data.ProviderName,
		Datacenter:   resDetails.Data.Datacenter,
		ServerType:   resDetails.Data.ServerType,
		Config: LoadBalancerConfig{
			HostHeader:             resDetails.Data.HostHeader,
			IsAccessLogsEnabled:    resDetails.Data.IsAccessLogsEnabled,
			IsForceHTTPSEnabled:    resDetails.Data.IsForceHTTPSEnabled,
			IPRateLimit:            resDetails.Data.IPRateLimit,
			IsIPRateLimitEnabled:   resDetails.Data.IsIPRateLimitEnabled,
			OutputCacheInSeconds:   resDetails.Data.OutputCacheInSeconds,
			IsStickySessionEnabled: resDetails.Data.IsStickySessionEnabled,
			IsProxyProtocolEnabled: resDetails.Data.IsProxyProtocolEnabled,
			SSLDomains:             resDetails.Data.SSLDomains,
			ForwardRules:           resDetails.Data.ForwardRules,
			OutputHeaders:          resDetails.Data.OutputHeaders,
			TargetServices:         resDetails.Data.TargetServices,
			RemoveResponseHeaders:  resDetails.Data.RemoveResponseHeaders,
		},
	}

	return &loadBalancer, nil
}

type CreateLoadBalancerRequest struct {
	ProjectID    string                          `json:"projectId"`
	ProviderName string                          `json:"providerName"`
	Datacenter   string                          `json:"datacenter"`
	ServerType   string                          `json:"serverType"`
	Config       CreateLoadBalancerRequestConfig `json:"loadBalancerPayload"`
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
