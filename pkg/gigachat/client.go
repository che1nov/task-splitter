package gigachat

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"task-splitter/config"
	"time"

	"github.com/google/uuid"
)

// Client представляет клиент для работы с GigaChat API
type Client struct {
	clientID string
	scope    string
	authKey  string
	client   *http.Client
}

// NewClient создает новый клиент GigaChat
func NewClient(cfg config.GigaChatConfig) *Client {
	// Создаем HTTP клиент с отключенной проверкой SSL для GigaChat
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		clientID: cfg.ClientID,
		scope:    cfg.Scope,
		authKey:  cfg.AuthKey,
		client: &http.Client{
			Transport: tr,
			Timeout:   60 * time.Second,
		},
	}
}

// SplitTask отправляет запрос на разбивку задачи в GigaChat
func (c *Client) SplitTask(ctx context.Context, prompt string) (string, error) {
	if c.authKey == "" {
		// Возвращаем заглушку для разработки
		return c.getMockResponse(), nil
	}

	// Сначала получаем access token
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	// Формируем запрос к GigaChat
	requestBody := map[string]interface{}{
		"model": "GigaChat",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Ты эксперт по планированию и разбивке задач. Помоги разбить задачу на логические подзадачи. Отвечай в формате JSON массива с объектами, содержащими поля: title, description, priority, estimated_time.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Разбей следующую задачу на подзадачи: %s", prompt),
			},
		},
		"max_tokens":  2000,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GigaChat API error: %s", string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}

	return content, nil
}

// getAccessToken получает access token для GigaChat API
func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	if c.authKey == "" {
		return "", fmt.Errorf("GigaChat auth key is not configured")
	}

	// Используем form-encoded данные для GigaChat API
	data := url.Values{}
	data.Set("scope", c.scope)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+c.authKey)
	req.Header.Set("RqUID", uuid.New().String())

	fmt.Printf("GigaChat token request - URL: %s, Data: %s, Auth: Basic %s, RqUID: %s\n", req.URL.String(), data.Encode(), c.authKey, req.Header.Get("RqUID"))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	fmt.Printf("GigaChat token request - Status: %d, Body: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GigaChat token API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal token response: %w", err)
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("invalid access token format")
	}

	return accessToken, nil
}

// getMockResponse возвращает заглушку для разработки
func (c *Client) getMockResponse() string {
	return `[
		{
			"title": "Изучить основы Docker",
			"description": "Понять что такое контейнеризация и зачем она нужна",
			"priority": "high",
			"estimated_time": 120
		},
		{
			"title": "Установить Docker",
			"description": "Скачать и установить Docker Desktop на локальную машину",
			"priority": "high",
			"estimated_time": 30
		},
		{
			"title": "Изучить Dockerfile",
			"description": "Научиться создавать образы с помощью Dockerfile",
			"priority": "medium",
			"estimated_time": 180
		},
		{
			"title": "Работа с Docker Compose",
			"description": "Изучить оркестрацию контейнеров с помощью Docker Compose",
			"priority": "medium",
			"estimated_time": 150
		},
		{
			"title": "Практические примеры",
			"description": "Создать несколько практических проектов с использованием Docker",
			"priority": "low",
			"estimated_time": 240
		}
	]`
}
