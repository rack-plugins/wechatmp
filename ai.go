package wechatmp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fimreal/goutils/ezap"
)

var LLM LLMInstance

type LLMInstance interface {
	Ask(question string) (answer string, err error)
	SetModelName(name string)
	SetModelEndpoint(endpoint string)
	SetModelPrompt(prompt string)
	SetSafetyMode(enabled bool)
}

// var MLM MLMInstance

// type MLMInstance interface {
// 	Ask(prompt string) (answer string, err error)
// 	Draw(prompt string) (answer string, err error)
// 	SetModelName(name string)
// 	SetModelEndpoint(endpoint string)
// 	SetModelPrompt(prompt string)
// 	SetSafetyMode(enabled bool)
// }

// Txt2img 向指定接口发送请求生成动漫风格图片，返回图片链接。
func Txt2Anime(prompt string) (string, error) {
	// coze 免费接口
	url := "https://epurs.com/txt2img"

	type DrawRequest struct {
		UserID         string `json:"user_id"`
		BotID          string `json:"bot_id"`
		Prompt         string `json:"prompt"`
		ConversationID string `json:"conversation_id,omitempty"`
	}

	// 创建请求体
	requestBody := DrawRequest{
		UserID: "1122333",
		BotID:  "7396854937780551734", // 测试默认动漫风格
		Prompt: prompt,
		// ConversationID: conversationID,
	}

	// 编码为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal draw request body: %w", err)
	}
	ezap.Debugf("request body: %s", string(jsonData))

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	type ResponseBody struct {
		Content string `json:"content"`
	}

	// 解析响应体
	var responseBody ResponseBody

	if err = json.Unmarshal(body, &responseBody); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// 从内容中提取图片链接
	imageURL := extractImageURL(responseBody.Content)

	ezap.Debugf("获取到 image URL: %s", imageURL)
	return imageURL, nil
}

func extractImageURL(content string) string {
	// 简单提取链接的逻辑，可以根据需要进行调整
	start := strings.Index(content, "http")
	end := strings.Index(content[start:], ")")
	if start == -1 || end == -1 {
		return ""
	}
	return content[start : start+end]
}
