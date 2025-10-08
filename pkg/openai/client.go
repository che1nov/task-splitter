package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"task-splitter/config"
	"time"
)

// Client представляет клиент для работы с OpenAI API
type Client struct {
	apiKey string
	model  string
	client *http.Client
}

// NewClient создает новый клиент OpenAI
func NewClient(cfg config.OpenAIConfig) *Client {
	return &Client{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// SplitTask отправляет запрос на разбивку задачи в OpenAI
func (c *Client) SplitTask(ctx context.Context, prompt string) (string, error) {
	if c.apiKey == "" {
		// Возвращаем заглушку для разработки
		return c.getMockResponse(), nil
	}

	requestBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Ты эксперт по планированию и разбивке задач. Помоги разбить задачу на логические подзадачи.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  2000,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

// getMockResponse возвращает заглушку для разработки
func (c *Client) getMockResponse() string {
	return `[
		{
			"title": "Анализ требований",
			"description": "Изучить и проанализировать все требования к проекту",
			"priority": "high",
			"estimated_time": 120
		},
		{
			"title": "Планирование архитектуры",
			"description": "Спроектировать архитектуру системы",
			"priority": "high",
			"estimated_time": 180
		},
		{
			"title": "Настройка окружения разработки",
			"description": "Подготовить рабочее окружение и необходимые инструменты",
			"priority": "medium",
			"estimated_time": 60
		},
		{
			"title": "Реализация основных компонентов",
			"description": "Разработать основные компоненты системы",
			"priority": "high",
			"estimated_time": 480
		},
		{
			"title": "Тестирование",
			"description": "Провести тестирование всех компонентов",
			"priority": "medium",
			"estimated_time": 240
		},
		{
			"title": "Документация",
			"description": "Написать техническую документацию",
			"priority": "low",
			"estimated_time": 120
		}
	]`
}
