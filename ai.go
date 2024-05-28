package wechatmp

import (
	"strings"
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
	SetSafetyMode(enabled bool)
}

func init() {
	go func() {
		// 确保启动加载变量后再执行
		time.Sleep(3 * time.Second)

		modelName := viper.GetString(ID + ".modelname")
		modelEndpoint := viper.GetString(ID + ".modelendpoint")
		modelPrompt := viper.GetString(ID + ".modelprompt")
		modelApiKey := viper.GetString(ID + ".modelapikey")
		safetymode := viper.GetBool(ID + ".safetymode")

		switch {
		case strings.HasPrefix(modelName, "gemini"):
			LLM = gemini.NewSession(modelApiKey)
		// case strings.HasPrefix(modelName, "anotherModel"):
		//  LLM = anotherModel.NewSession(modelApiKey)
		default:
			LLM = gemini.NewSession(modelApiKey)
		}

		// 设置模型名称、提示语、安全模式
		LLM.SetModelPrompt(modelPrompt)
		LLM.SetModelName(modelName)
		LLM.SetModelEndpoint(modelEndpoint)
		LLM.SetSafetyMode(safetymode)
	}()
}
