package wechatmp

var LLM LLMInstance

type LLMInstance interface {
	Ask(question string) (answer string, err error)
	SetModelName(name string)
	SetModelEndpoint(endpoint string)
	SetModelPrompt(prompt string)
	SetSafetyMode(enabled bool)
}
