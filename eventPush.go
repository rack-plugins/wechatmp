package wechatmp

// 处理消息事件
func (m *WechatmpMessage) HandleEvent() (reply string, err error) {
	switch m.Event {
	case "subscribe":
		m.subscribeEvent()
	case "ubsubscribe":
		// 删掉取消订阅的用户会话
		MsgContext.Remove(m.FromUserName)
	case "CLICK":
	default:
	}
	return
}

// 处理关注/取消关注公众号事件
func (m *WechatmpMessage) subscribeEvent() (reply string, err error) {
	reply = "感谢您的关注, 请回复任意内容以开始对话。⚠️由于个人公众号限制原因，如果遇到过长内容超过10s未回复，请在15分钟内发送【重试】获取上一条回复。"
	return
}
