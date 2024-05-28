package wechatmp

import (
	"github.com/rack-plugins/wechatmp/gemini"
	"github.com/spf13/viper"
)

func askGemini(question string) (answer string, err error) {
	// 创建 llm 会话
	llm := gemini.NewSession(viper.GetString(ID + ".geminiapikey"))
	llm.SetModelPrompt("你是一个没有名字的人工智能助手,回答问题时尽量口语化,不要使用markdown文本标记.")

	return llm.Ask(question)
}
