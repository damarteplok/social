package zeebe

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

func NewTokenManager(clientID, clientSecret, authURL string) *TokenManager {
	return &TokenManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		authURL:      authURL,
	}
}

func (t *TokenManager) isTokenExpired() bool {
	return time.Now().After(t.expiry)
}

func (t *TokenManager) refreshTokenFunc() error {
	formData := url.Values{}
	formData.Set("client_id", t.clientID)
	formData.Set("client_secret", t.clientSecret)
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", t.refreshToken)

	resp, err := http.PostForm(t.authURL, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to refresh token")
	}

	var data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	t.authToken = data.AccessToken
	t.refreshToken = data.RefreshToken
	t.expiry = time.Now().Add(time.Duration(data.ExpiresIn) * time.Second)

	return nil
}

func (t *TokenManager) loginFromApp() error {
	formData := url.Values{}
	formData.Set("client_id", t.clientID)
	formData.Set("client_secret", t.clientSecret)
	formData.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(t.authURL, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to login")
	}

	var data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	t.authToken = data.AccessToken
	t.refreshToken = data.RefreshToken
	t.expiry = time.Now().Add(time.Duration(data.ExpiresIn) * time.Second)

	return nil
}

func (t *TokenManager) GetAuthToken(ctx context.Context) (string, error) {
	if t.authToken != "" {
		if t.isTokenExpired() {
			err := t.refreshTokenFunc()
			if err != nil {
				return "", err
			}
		}
	} else {
		err := t.loginFromApp()
		if err != nil {
			return "", err
		}
	}
	return t.authToken, nil
}
