package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/leetsecure/qryptic-client-cli/internal/models"
)

type QrypticClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewLeetClient creates a new instance of LeetClient.
func NewQrypticClient(baseURL, authToken string) *QrypticClient {
	return &QrypticClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		AuthToken: authToken,
	}
}

func (c *QrypticClient) EmailPasswordLogin(req models.EmailPasswordLoginRequest) (int, *models.AuthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/login", c.BaseURL)

	statusCode, respBody, err := c.doRequest(http.MethodPost, url, req)
	if err != nil {
		return 0, nil, err
	}

	var response models.AuthResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return 0, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return statusCode, &response, nil
}

func (c *QrypticClient) ControllerHealthCheck() (int, *models.HealthCheckResponse, error) {
	url := fmt.Sprintf("%s/api/v1/health", c.BaseURL)

	statusCode, respBody, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}

	var response models.HealthCheckResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return 0, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return statusCode, &response, nil
}

func (c *QrypticClient) ListAccessibleGateways() (int, *([]models.GatewayResponse), error) {
	url := fmt.Sprintf("%s/api/v1/gateway/list", c.BaseURL)

	statusCode, respBody, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}

	var response []models.GatewayResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return 0, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return statusCode, &response, nil
}

func (c *QrypticClient) GetGatewayClient(uuid string) (int, *(models.WGClientConfig), error) {
	url := fmt.Sprintf("%s/api/v1/gateway/%s/client", c.BaseURL, uuid)

	statusCode, respBody, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}

	var clientConfigResponse models.WGClientConfig
	if err := json.Unmarshal(respBody, &clientConfigResponse); err != nil {
		return 0, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return statusCode, &clientConfigResponse, nil
}

func (c *QrypticClient) GetWebSSOToken(codeVerifier, codeChallenge string) (int, *models.AuthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/google/web/sso/token?code_verifier=%s&code_challenge=%s", c.BaseURL, codeVerifier, codeChallenge)

	statusCode, respBody, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}
	var response models.AuthResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return 0, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return statusCode, &response, nil
}

func (c *QrypticClient) doRequest(method, url string, body interface{}) (int, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// if resp.StatusCode >= 400 {
	// 	var errResp models.ErrorResponse
	// 	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
	// 		return 0, nil, fmt.Errorf("failed to parse error response: %w", err)
	// 	}
	// 	return resp.StatusCode, errors.New(errResp.Message)
	// }
	dataBytes, err := io.ReadAll(resp.Body)
	return resp.StatusCode, dataBytes, err
}
