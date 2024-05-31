package wechatmp

import (
	"context"
	"time"

	"github.com/fimreal/goutils/ezap"
)

// 处理文本消息
func (m *WechatmpMessage) HandleText() (reply string, err error) {
	// 被动回复消息
	return m.ReplyTextMessage()
}

// 构造回复消息 - 文本消息
func (message *WechatmpMessage) ReplyTextMessage() (reply string, err error) {
	// 处理重试消息
	if message.Content == "重试" {
		reply = "未获取到历史消息，请尝试重新发送问题吧"
		if oldMsgId := MsgContext.GetOrSet(message.FromUserName, int64(0)).(int64); oldMsgId != int64(0) {
			out := MsgContext.Get(oldMsgId).(string)
			if out != "" {
				reply = out
				ezap.Infof("获取到重试回复的消息: msgId: %d", oldMsgId)
			}
			ezap.Debugf("未获取到重试回复的消息: msgId: %d", oldMsgId)
		}
		return
	}

	// 处理 msgId 相同的重复请求，返回历史消息
	if MsgContext.GetOrSet(message.FromUserName, int64(0)).(int64) == message.MsgId {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		replyChan := make(chan string, 1)
		go func() {
			for {
				reply := MsgContext.Get(message.MsgId)
				if reply != "" {
					replyChan <- reply.(string)
					return
				}
				time.Sleep(100 * time.Microsecond)
			}
		}()

		select {
		case reply = <-replyChan:
		case <-ctx.Done():
		}

		ezap.Infof("获取到超时回复的消息: msgId: %d", message.MsgId)
		return
	}

	// 正常流程，询问 llm
	reply, err = LLM.Ask(message.Content)
	if err == nil {
		// 保存新消息到用户对话上下文
		MsgContext.Set(message.FromUserName, message.MsgId)
		MsgContext.Set(message.MsgId, reply)
		// 15 分钟后删除历史消息，减少内存占用
		go func() {
			time.Sleep(15 * time.Minute)
			MsgContext.Remove(message.MsgId)
		}()
	}
	return
}
