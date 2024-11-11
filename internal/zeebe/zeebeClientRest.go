package zeebe

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/damarteplok/social/internal/store"
)

func NewZeebeClientRest(clientID, clientSecret, authServerURL, zeebeAddr string) (*ZeebeClientRest, error) {
	tokenManager := NewTokenManager(clientID, clientSecret, authServerURL)

	// Fetch initial token
	_, err := tokenManager.GetAuthToken(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get initial auth token: %w", err)
	}

	// Create an HTTP client
	httpClient := &http.Client{
		Transport: &http.Transport{},
	}

	return &ZeebeClientRest{
		httpClient:   httpClient,
		tokenManager: tokenManager,
		zeebeAddr:    zeebeAddr,
	}, nil
}

func (z *ZeebeClientRest) Close() error {
	return nil
}

func (z *ZeebeClientRest) SendRequest(ctx context.Context, method, endpoint string, body io.Reader) ([]byte, error) {
	token, err := z.tokenManager.GetAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, store.ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}
