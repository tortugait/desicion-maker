package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cohesion-org/deepseek-go"
)

type (
	DeepSeekConf struct {
		BaseURL string
		APIKey  string
	}

	deepSeekSrv struct {
		conf DeepSeekConf
	}
)

func NewDeepSeekSrv(conf DeepSeekConf) deepSeekSrv {
	return deepSeekSrv{conf: conf}
}

func (s deepSeekSrv) GenerateReason(ctx context.Context, statement string, ok bool) (string, error) {
	client := deepseek.NewClient(s.conf.APIKey)

	systemMsg := "Create humorous explanations for boolean answers. " +
		"Keep responses under 150 characters. Use the same language as the input. " +
		"Return only the reason itself."

	userPrompt := fmt.Sprintf("Statement: %s\nVeracity: %t\nGenerate reason:", statement, ok)

	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: systemMsg},
			{Role: deepseek.ChatMessageRoleUser, Content: userPrompt},
		},
	}

	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to create completion: %w", err) // Proper error handling
	}

	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}
