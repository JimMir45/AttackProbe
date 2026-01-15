package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ClientConfig struct {
	Endpoint     string
	APIKey       string
	Model        string
	Timeout      int    // ms
	ExtraHeaders string // JSON格式
}

type Client interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	Test() error
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Content      string `json:"content"`
	Model        string `json:"model"`
	FinishReason string `json:"finish_reason"`
	Usage        Usage  `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAIClient OpenAI兼容客户端
type OpenAIClient struct {
	endpoint     string
	apiKey       string
	model        string
	httpClient   *http.Client
	extraHeaders map[string]string
}

func NewOpenAIClient(config *ClientConfig) *OpenAIClient {
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30000
	}

	client := &OpenAIClient{
		endpoint: strings.TrimSuffix(config.Endpoint, "/"),
		apiKey:   config.APIKey,
		model:    config.Model,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
		extraHeaders: make(map[string]string),
	}

	// 解析额外请求头
	if config.ExtraHeaders != "" {
		json.Unmarshal([]byte(config.ExtraHeaders), &client.extraHeaders)
	}

	return client
}

func (c *OpenAIClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = c.model
	}

	body := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	endpoint := c.endpoint
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint = endpoint + "/chat/completions"
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	for k, v := range c.extraHeaders {
		httpReq.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败: %d - %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
		} `json:"choices"`
		Usage Usage  `json:"usage"`
		Model string `json:"model"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Error != nil {
		return nil, errors.New(result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("LLM无响应")
	}

	return &ChatResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		FinishReason: result.Choices[0].FinishReason,
		Usage:        result.Usage,
	}, nil
}

func (c *OpenAIClient) Test() error {
	// 使用客户端配置的超时时间，而不是硬编码10秒
	ctx, cancel := context.WithTimeout(context.Background(), c.httpClient.Timeout)
	defer cancel()

	_, err := c.Chat(ctx, &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "Hi"},
		},
		MaxTokens: 5,
	})

	return err
}
