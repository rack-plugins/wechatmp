package wechatmp

import "github.com/spf13/viper"

// 处理消息事件
func (m *WechatmpMessage) HandleEvent() (reply string, err error) {
	switch m.Event {
	case "subscribe":
		return m.subscribeEvent()
	case "ubsubscribe":
		// 删掉取消订阅的用户会话
		MsgContext.Remove(m.FromUserName)
	case "CLICK":
	default:
	}
	return
}

// 处理关注公众号事件，返回欢迎语
func (m *WechatmpMessage) subscribeEvent() (reply string, err error) {
	reply = viper.GetString(ID + ".subscribeMessage")
	return
}
