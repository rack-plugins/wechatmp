package wechatmp

import (
	"time"

	"github.com/fimreal/rack/module"
	"github.com/spf13/cobra"
)

const (
	ID            = "wechatmp"
	Comment       = "[module] wechatmp api service"
	RoutePrefix   = "/wx"
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

	serveCmd.Flags().String(ID+".geminiapikey", "", "Gemini_API_Token")
}

func init() {
	// 待启动加载参数后再执行
	go func() {
		time.Sleep(3 * time.Second)
		refreshAccessToken()
	}()
}
