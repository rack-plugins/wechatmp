package wechatmp

import (
	"strings"
	"time"

	"github.com/fimreal/rack/module"
	"github.com/rack-plugins/wechatmp/gemini"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ID            = "wechatmp"
	Comment       = "[module] wechatmp api service"
	RoutePrefix   = "/"
	DefaultEnable = false
)

var Module = module.Module{
	ID:      ID,
	Comment: Comment,
	// gin route
	RouteFunc:   AddRoute,
	RoutePrefix: RoutePrefix,
	// cobra flag
	FlagFunc: ServeFlag,
	// crond
	// CrondFunc: map[string]func(){"[wechatmp.refreshAccessToken] @every 1h30m": RefreshAccessToken},
}

func ServeFlag(serveCmd *cobra.Command) {
	serveCmd.Flags().Bool(ID, DefaultEnable, Comment)

	serveCmd.Flags().String(ID+".appid", "", "公众号开发者 ID (AppID)")
	serveCmd.Flags().String(ID+".appsecret", "", "公众号开发者密码 (AppSecret)")
	serveCmd.Flags().String(ID+".token", "", "公众号令牌 (Token)")

	serveCmd.Flags().String(ID+".subscribeMessage", "感谢您的关注", "订阅公众号时返回信息")

	serveCmd.Flags().String(ID+".modelapikey", "", "Gemini_API_Token")
	serveCmd.Flags().String(ID+".modelname", "gemini-1.5-pro-latest", "默认使用 gemini-1.5-pro-latest，目前仅支持免费的 gemini")
	serveCmd.Flags().String(ID+".modelprompt", "你是一个没有名字的人工智能助手,回答问题时尽量口语化,不要使用markdown文本标记.", "设置提示词")
	serveCmd.Flags().String(ID+".modelendpoint", "generativelanguage.googleapis.com", "设置模型地址")
	serveCmd.Flags().Bool(ID+".safetymode", false, "开启安全模式")
}

func init() {
	if !viper.GetBool(ID) && !viper.GetBool("allservices") {
		return
	}

	// 待启动加载参数后再执行
	go func() {
		time.Sleep(3 * time.Second)
		refreshAccessToken()
	}()
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
