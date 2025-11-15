package domain

import (
	"strings"
	"time"
)

// Template представляет шаблон для разбивки задач
type Template struct {
	ID          uint
	Name        string
	Description string
	Category    string
	Prompt      string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTemplate создает новый шаблон с валидацией
func NewTemplate(name, description, category, prompt string) (Template, error) {
	// Валидация названия
	name = strings.TrimSpace(name)
	if name == "" {
		return Template{}, ErrEmptyTemplateName
	}

	// Валидация промпта
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return Template{}, ErrEmptyTemplatePrompt
	}

	return Template{
		Name:        name,
		Description: strings.TrimSpace(description),
		Category:    strings.TrimSpace(category),
		Prompt:      prompt,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// Update обновляет шаблон
func (t *Template) Update(name, description, category, prompt string) error {
	// Валидация и обновление названия
	if name != "" {
		name = strings.TrimSpace(name)
		if name == "" {
			return ErrEmptyTemplateName
		}
		t.Name = name
	}

	// Обновление описания
	if description != "" {
		t.Description = strings.TrimSpace(description)
	}

	// Обновление категории
	if category != "" {
		t.Category = strings.TrimSpace(category)
	}

	// Валидация и обновление промпта
	if prompt != "" {
		prompt = strings.TrimSpace(prompt)
		if prompt == "" {
			return ErrEmptyTemplatePrompt
		}
		t.Prompt = prompt
	}

	t.UpdatedAt = time.Now()
	return nil
}

// Activate активирует шаблон
func (t *Template) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

// Deactivate деактивирует шаблон
func (t *Template) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

