package wechatmp

import (
	"time"

	"github.com/spf13/viper"
)

// 处理消息事件
func (m *WechatmpMessage) HandleEvent() (reply *WechatmpMessage, err error) {
	// 创建回复消息模板
	reply = &WechatmpMessage{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgId:        m.MsgId,
		MsgType:      "text",
	}
	switch *m.Event {
	case "subscribe":
		msg := m.subscribeMessage()
		reply.Content = &msg
		return
	case "ubsubscribe":
		// 删掉取消订阅的用户会话
		MsgContext.Remove(m.FromUserName)
	case "CLICK":
	default:
	}
	return
}

// 处理关注公众号事件，返回欢迎语
func (m *WechatmpMessage) subscribeMessage() (conteng string) {
	subscribeMessage := viper.GetString(ID + ".subscribeMessage")
	if subscribeMessage == "" {
		subscribeMessage = "欢迎关注！"
	}
	m.Content = &subscribeMessage
	return
}
