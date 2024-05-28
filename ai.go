package wechatmp

import (
	"time"

	"github.com/rack-plugins/wechatmp/gemini"
	"github.com/spf13/viper"
)

var LLM LLMInstance

type LLMInstance interface {
	Ask(question string) (answer string, err error)
	SetModelName(name string)
	SetModelEndpoint(endpoint string)
	SetModelPrompt(prompt string)
	SetSatetyMode(enabled bool)
}

func init() {
	go func() {
		// 待启动加载参数后再执行
		time.Sleep(3 * time.Second)
		// 创建 llm 会话
		LLM = gemini.NewSession(viper.GetString(ID + ".geminiapikey"))
		LLM.SetModelPrompt(viper.GetString(ID + ".modelprompt"))
		LLM.SetModelName(viper.GetString(ID + ".modelname"))
		LLM.SetModelEndpoint(viper.GetString(ID + ".modelendpoint"))
		LLM.SetSatetyMode(viper.GetBool(ID + ".safetymode"))
	}()
}
