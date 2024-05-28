package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/fimreal/goutils/ezap"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Session struct {
	Name       string
	Token      string
	Endpoint   string
	Prompt     genai.Part
	SafetyMode bool // default is false
}

func NewSession(token string) *Session {
	return &Session{
		// https://cloud.google.com/vertex-ai/generative-ai/docs/learn/model-versioning?hl=zh-cn#stable-versions-available for more information
		Name:     "gemini-1.0-pro",
		Token:    token,
		Endpoint: "generativelanguage.googleapis.com",
		Prompt:   nil,
	}
}

func (s *Session) SetModelName(name string) {
	s.Name = name
	ezap.Info("Set Model Name: ", s.Name)
}

func (s *Session) SetModelEndpoint(endpoint string) {
	s.Endpoint = endpoint
	ezap.Info("Set Model Endpoint: ", s.Endpoint)
}

func (s *Session) SetModelPrompt(prompt string) {
	s.Prompt = genai.Text(prompt)
	ezap.Info("Set Model Prompt: ", s.Prompt)
}

func (s *Session) SetSafetyMode(enabled bool) {
	s.SafetyMode = enabled
	ezap.Info("Set Safety Mode: ", s.SafetyMode)
}

func (s *Session) Ask(question string) (answer string, err error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.Token), option.WithEndpoint(s.Endpoint))
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Create a Generative Model
	model := client.GenerativeModel(s.Name)
	// model.SetTemperature(0.9)
	// model.SetTopP(0.5)
	// model.SetTopK(20)
	// model.SetMaxOutputTokens(100)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{s.Prompt},
	}
	// satety mode
	if s.SafetyMode {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockMediumAndAbove,
			},
		}
	}

	ezap.Info("Ask Gemini: ", question)
	resp, err := model.GenerateContent(ctx, genai.Text(question))
	if err != nil {
		return "", err
	}
	ctx.Done()
	return readAllFrom(resp), err
}

func readAllFrom(resp *genai.GenerateContentResponse) string {
	var parts strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				ezap.Debug("Gemini Say: ", part)
				parts.WriteString(fmt.Sprintf("%s", part))
			}
		}
	}
	return parts.String()
}
