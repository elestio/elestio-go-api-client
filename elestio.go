package elestio

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	BaseURLV1 = "https://api.elest.io"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Email      string
	ApiKey     string
	jwt        string

	Project      *ProjectHandler
	Service      *ServiceHandler
	LoadBalancer *LoadBalancerHandler
}

func NewClient(email, apiKey string) (*Client, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	if apiKey == "" {
		return nil, errors.New("api key is required")
	}

	client := Client{
		BaseURL:    BaseURLV1,
		HTTPClient: &http.Client{},
		Email:      email,
		ApiKey:     apiKey,
	}

	if err := client.signIn(); err != nil {
		return nil, fmt.Errorf("failed to sign in: %s", err)
	}

	client.init()

	return &client, nil
}

func NewUnsignedClient() *Client {
	client := Client{
		BaseURL:    BaseURLV1,
		HTTPClient: &http.Client{},
	}

	client.init()

	return &client
}

// init sets up all the handlers.
func (c *Client) init() {
	c.Project = &ProjectHandler{client: c}
	c.Service = &ServiceHandler{client: c}
	c.LoadBalancer = &LoadBalancerHandler{client: c}
}
