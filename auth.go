package elestio

import (
	"fmt"
)

type (
	authRequest struct {
		Email  string `json:"email"`
		ApiKey string `json:"token"`
	}

	authResponse struct {
		APIResponse
		JWT string `json:"jwt"`
	}
)

func (c *Client) signIn() error {
	bts, err := c.sendPostRequest(fmt.Sprintf("%s/api/auth/checkAPIToken", c.BaseURL), authRequest{c.Email, c.ApiKey})
	if err != nil {
		return err
	}

	var r authResponse
	if err := checkAPIResponse(bts, &r); err != nil {
		return err
	}

	c.jwt = r.JWT

	return nil
}
