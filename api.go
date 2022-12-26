package elestio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	// APIResponse represents a response returned by Elestio.
	APIResponse struct {
		Status  string `json:"status,omitempty"`
		Message string `json:"message,omitempty"`
	}
)

func checkAPIResponse(bts []byte, r any) error {
	if r == nil {
		r = new(APIResponse)
	}

	buffer := bytes.NewBuffer(bts)
	dec := json.NewDecoder(buffer)
	if err := dec.Decode(&r); err != nil {
		return fmt.Errorf("cannot unmarshal JSON `%s`, error: %w", bts, err)
	}

	return nil
}

func (c *Client) sendGetRequest(endpoint string, req interface{}) ([]byte, error) {
	return c.sendRequest("GET", endpoint, req)
}

func (c *Client) sendPutRequest(endpoint string, req interface{}) ([]byte, error) {
	return c.sendRequest("PUT", endpoint, req)
}

func (c *Client) sendPostRequest(endpoint string, req interface{}) ([]byte, error) {
	return c.sendRequest("POST", endpoint, req)
}

func (c *Client) sendDeleteRequest(endpoint string, req interface{}) ([]byte, error) {
	return c.sendRequest("DELETE", endpoint, req)
}

func (c *Client) sendRequest(method string, url string, body any) ([]byte, error) {
	var bts []byte
	if body != nil {
		var err error
		bts, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	retryCount := 2
	for {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(bts))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.jwt))

		// Temporary fix waiting api handle jwt in authoization header
		query := req.URL.Query()
		query.Set("jwt", c.jwt)
		req.URL.RawQuery = query.Encode()

		rsp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer func() {
			err := rsp.Body.Close()
			if err != nil {
				log.Println("Cannot close response body: %w", err)
			}
		}()

		responseBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		// Retry in case of timeout or error for GET requests
		if (rsp.StatusCode == 408 || rsp.StatusCode >= 500) && retryCount > 0 && method == "GET" {
			retryCount--
			continue
		} else if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
			return nil, fmt.Errorf("request failed with status code %d: %s", rsp.StatusCode, string(responseBody))
		}

		return responseBody, nil
	}
}
