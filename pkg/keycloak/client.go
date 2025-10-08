package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"task-splitter/config"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client представляет клиент для работы с Keycloak
type Client struct {
	config       config.KeycloakConfig
	httpClient   *http.Client
	oauth2Config *oauth2.Config
}

// UserInfo представляет информацию о пользователе из Keycloak
type UserInfo struct {
	Sub               string   `json:"sub"`
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	RealmRoles        []string `json:"realm_roles"`
	Groups            []string `json:"groups"`
}

// TokenInfo представляет информацию о токене
type TokenInfo struct {
	Active     bool     `json:"active"`
	Sub        string   `json:"sub"`
	Email      string   `json:"email"`
	Username   string   `json:"preferred_username"`
	Exp        int64    `json:"exp"`
	RealmRoles []string `json:"realm_roles"`
	Groups     []string `json:"groups"`
}

// NewClient создает новый клиент Keycloak
func NewClient(cfg config.KeycloakConfig) *Client {
	// Проверяем базовую конфигурацию
	if cfg.URL == "" || cfg.Realm == "" || cfg.ClientID == "" {
		fmt.Println("⚠️  Keycloak configuration incomplete, skipping connection test")
		return nil
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", cfg.URL, cfg.Realm),
			TokenURL: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", cfg.URL, cfg.Realm),
		},
	}

	client := &Client{
		config:       cfg,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		oauth2Config: oauth2Config,
	}

	// Проверяем подключение к Keycloak
	wellKnownURL := fmt.Sprintf("%s/realms/%s/.well-known/openid_configuration", cfg.URL, cfg.Realm)
	resp, err := client.httpClient.Get(wellKnownURL)
	if err != nil {
		fmt.Printf("⚠️  Failed to connect to Keycloak: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️  Keycloak server returned status %d\n", resp.StatusCode)
		return nil
	}

	fmt.Println("✅ Keycloak connected successfully")
	return client
}

// ValidateToken проверяет валидность токена
func (c *Client) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	introspectURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
		c.config.URL, c.config.Realm)

	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed: %s", string(body))
	}

	var tokenInfo TokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !tokenInfo.Active {
		return nil, fmt.Errorf("token is not active")
	}

	return &tokenInfo, nil
}

// GetUserInfo получает информацию о пользователе
func (c *Client) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo",
		c.config.URL, c.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &userInfo, nil
}

// GetAdminToken получает токен администратора для управления пользователями
func (c *Client) GetAdminToken(ctx context.Context) (string, error) {
	config := &clientcredentials.Config{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		TokenURL:     fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.config.URL, c.config.Realm),
		Scopes:       []string{"openid"},
	}

	token, err := config.Token(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get admin token: %w", err)
	}

	return token.AccessToken, nil
}

// CreateUser создает пользователя в Keycloak
func (c *Client) CreateUser(ctx context.Context, username, email, firstName, lastName, password string) error {
	adminToken, err := c.GetAdminToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin token: %w", err)
	}

	createUserURL := fmt.Sprintf("%s/admin/realms/%s/users", c.config.URL, c.config.Realm)

	userData := map[string]interface{}{
		"username":  username,
		"email":     email,
		"firstName": firstName,
		"lastName":  lastName,
		"enabled":   true,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", createUserURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user: %s", string(body))
	}

	return nil
}

// AssignRoleToUser назначает роль пользователю
func (c *Client) AssignRoleToUser(ctx context.Context, userID, roleName string) error {
	adminToken, err := c.GetAdminToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin token: %w", err)
	}

	// Получаем роль
	role, err := c.getRole(ctx, adminToken, roleName)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Назначаем роль пользователю
	assignRoleURL := fmt.Sprintf("%s/admin/realms/%s/users/%s/role-mappings/realm",
		c.config.URL, c.config.Realm, userID)

	roleData := []map[string]interface{}{role}

	jsonData, err := json.Marshal(roleData)
	if err != nil {
		return fmt.Errorf("failed to marshal role data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", assignRoleURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to assign role: %s", string(body))
	}

	return nil
}

// getRole получает роль по имени
func (c *Client) getRole(ctx context.Context, adminToken, roleName string) (map[string]interface{}, error) {
	getRoleURL := fmt.Sprintf("%s/admin/realms/%s/roles/%s", c.config.URL, c.config.Realm, roleName)

	req, err := http.NewRequestWithContext(ctx, "GET", getRoleURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get role: %s", string(body))
	}

	var role map[string]interface{}
	if err := json.Unmarshal(body, &role); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return role, nil
}
