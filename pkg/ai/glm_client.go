package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GLMClient 智谱GLM大模型客户端
type GLMClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	Endpoint   string
	HTTPClient *http.Client
}

// GLMRequest GLM请求结构（与OpenAI兼容）
type GLMRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// GLMResponse GLM响应结构（与OpenAI兼容）
type GLMResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GLMError GLM错误结构
type GLMError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// NewGLMClient 创建GLM客户端
// baseURL: 智谱API基础地址，默认 https://open.bigmodel.cn
// apiKey: 智谱API Key
// model: 模型名称，如 glm-4, glm-4-flash, glm-3-turbo
func NewGLMClient(baseURL, apiKey, model, endpoint string) *GLMClient {
	if baseURL == "" {
		baseURL = "https://open.bigmodel.cn"
	}
	if endpoint == "" {
		endpoint = "/api/paas/v4/chat/completions"
	}
	if model == "" {
		model = "glm-4"
	}

	return &GLMClient{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Model:    model,
		Endpoint: endpoint,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Minute,
		},
	}
}

// GenerateText generates text using GLM API (implements AIClient interface)
func (c *GLMClient) GenerateText(prompt string, systemPrompt string, options ...func(*ChatCompletionRequest)) (string, error) {
	messages := []ChatMessage{}

	if systemPrompt != "" {
		messages = append(messages, ChatMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: prompt,
	})

	req := &GLMRequest{
		Model:    c.Model,
		Messages: messages,
	}

	Temp := float64(0)
	MaxT := int(0)
	TopP_val := float64(0)

	for _, option := range options {
		tempReq := &ChatCompletionRequest{}
		option(tempReq)
		if tempReq.Temperature > 0 {
			Temp = tempReq.Temperature
		}
		if tempReq.MaxTokens != nil {
			MaxT = *tempReq.MaxTokens
		}
		if tempReq.TopP > 0 {
			TopP_val = tempReq.TopP
		}
	}

	if Temp > 0 {
		req.Temperature = Temp
	}
	if MaxT > 0 {
		req.MaxTokens = &MaxT
	}
	if TopP_val > 0 {
		req.TopP = TopP_val
	}

	return c.sendRequest(req)
}

// sendRequest 发送请求
func (c *GLMClient) sendRequest(req *GLMRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("GLM: Failed to marshal request: %v\n", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.BaseURL + c.Endpoint

	// 打印请求信息（隐藏API Key）
	safeURL := url
	if strings.Contains(url, c.APIKey) {
		safeURL = strings.Replace(url, c.APIKey, "***", 1)
	}
	fmt.Printf("GLM: Sending request to: %s\n", safeURL)
	fmt.Printf("GLM: Model=%s\n", c.Model)

	requestPreview := string(jsonData)
	if len(jsonData) > 300 {
		requestPreview = string(jsonData[:300]) + "..."
	}
	fmt.Printf("GLM: Request body: %s\n", requestPreview)

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("GLM: Failed to create request: %v\n", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	fmt.Printf("GLM: Executing HTTP request...\n")
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		fmt.Printf("GLM: HTTP request failed: %v\n", err)
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("GLM: Received response with status: %d\n", resp.StatusCode)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("GLM: Failed to read response body: %v\n", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("GLM: API error (status %d): %s\n", resp.StatusCode, string(body))

		// 尝试解析错误信息
		var errResp struct {
			Error GLMError `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
			return "", fmt.Errorf("GLM API error: %s (code: %s)", errResp.Error.Message, errResp.Error.Code)
		}
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// 打印响应体（调试用）
	bodyPreview := string(body)
	if len(body) > 500 {
		bodyPreview = string(body[:500]) + "..."
	}
	fmt.Printf("GLM: Response body: %s\n", bodyPreview)

	// 解析响应
	var glmResp GLMResponse
	if err := json.Unmarshal(body, &glmResp); err != nil {
		errorPreview := string(body)
		if len(body) > 200 {
			errorPreview = string(body[:200])
		}
		fmt.Printf("GLM: Failed to parse response: %v\n", err)
		return "", fmt.Errorf("failed to unmarshal response: %w, body preview: %s", err, errorPreview)
	}

	fmt.Printf("GLM: Successfully parsed response, choices count: %d\n", len(glmResp.Choices))

	// 检查choices
	if len(glmResp.Choices) == 0 {
		fmt.Printf("GLM: No choices in response\n")
		return "", fmt.Errorf("no choices in response")
	}

	// 获取内容
	content := glmResp.Choices[0].Message.Content
	finishReason := glmResp.Choices[0].FinishReason

	fmt.Printf("GLM: finish_reason=%s, content_length=%d\n", finishReason, len(content))

	// 检查finish_reason
	if finishReason == "content_filter" {
		return "", fmt.Errorf("AI内容被安全过滤器拦截，可能因为：\n1. 请求内容触发了安全策略\n2. 生成的内容包含敏感信息\n3. 建议：调整输入内容或联系API提供商调整过滤策略")
	}

	if len(content) == 0 && finishReason != "stop" {
		return "", fmt.Errorf("AI返回内容为空 (finish_reason: %s)", finishReason)
	}

	fmt.Printf("GLM: Generated text: %s\n", content)
	return content, nil
}

// GenerateImage 生成图片（GLM暂不支持图像生成，返回错误）
func (c *GLMClient) GenerateImage(prompt string, size string, n int) ([]string, error) {
	return nil, fmt.Errorf("GLM client does not support image generation, use dedicated image API (CogView)")
}

// TestConnection 测试连接
func (c *GLMClient) TestConnection() error {
	fmt.Printf("GLM: TestConnection called with BaseURL=%s, Model=%s, Endpoint=%s\n", c.BaseURL, c.Model, c.Endpoint)

	_, err := c.GenerateText("你好", "你是一个友好的AI助手")
	if err != nil {
		fmt.Printf("GLM: TestConnection failed: %v\n", err)
	} else {
		fmt.Printf("GLM: TestConnection succeeded\n")
	}
	return err
}

// GLM支持的模型列表
var GLMModels = map[string]string{
	"glm-4":        "GLM-4: 智谱最新旗舰模型，支持超长上下文",
	"glm-4-flash":  "GLM-4-Flash: 免费版本，速度快，支持超长上下文",
	"glm-4-plus":   "GLM-4-Plus: GLM-4升级版，性能更强",
	"glm-3-turbo":  "GLM-3-Turbo: 性价比高，适合日常使用",
	"glm-4v":       "GLM-4V: 视觉理解模型，支持图像输入",
	"glm-4v-flash": "GLM-4V-Flash: 免费视觉理解模型",
	"cogview-3":    "CogView-3: 图像生成模型",
}

// GetAvailableModels 获取可用模型列表
func GetAvailableModels() map[string]string {
	return GLMModels
}
