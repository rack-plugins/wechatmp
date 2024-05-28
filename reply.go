package wechatmp

import (
	"context"
	"net/http"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/container/gmap"
	"github.com/spf13/viper"
)

var (
	MsgContext = gmap.New()
)

// 接收公众号消息
// https://mp.weixin.qq.com/debug/cgi-bin/apiinfo?t=index&type=%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8F%9C%E5%8D%95&form=%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8F%9C%E5%8D%95%E5%88%9B%E5%BB%BA%E6%8E%A5%E5%8F%A3%20/menu/creat
func ReplyWechatmpMessage(c *gin.Context) {
	// 获取微信加密签名
	var sig Signature
	sig.Signature = c.Query("signature")
	sig.Timestamp = c.Query("timestamp")
	sig.Nonce = c.Query("nonce")

	token := viper.GetString(ID + ".token")
	if !sig.checkSignature(token) {
		ezap.Error("签名验证失败")
		c.XML(http.StatusBadRequest, gin.H{"error": "签名验证失败"})
		return
	}
	ezap.Debug("签名验证成功")

	// 解析接收到的信息
	var message WechatmpMessage
	err := c.ShouldBind(&message)
	if err != nil {
		c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ezap.Debugf("%+v", message)

	// 生成回复消息
	reply := replyMessage(message)
	ezap.Debugf("回复消息: %+v", reply)

	c.XML(http.StatusOK, reply)
}

func replyMessage(message WechatmpMessage) *WechatmpMessage {
	// 创建回复消息模板
	reply := &WechatmpMessage{
		ToUserName:   message.FromUserName,
		FromUserName: message.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgId:        message.MsgId,
		MsgType:      "text",
	}

	// 根据消息类型回复
	switch message.MsgType {
	case "text":
		out, err := replyTextMessage(message)
		if err != nil {
			ezap.Error("构造文本消息回复时出错: ", err.Error())
			out = "出现错误，请联系管理员，稍后再来试一试吧。"
		}
		// 回复内容使用 CDATA 包裹, 解决换行问题
		reply.Content = "<![CDATA[" + out + "]]>"
	default:
		reply.Content = "抱歉，暂时无法处理此类型消息"
	}

	return reply
}

// 构造回复消息 - 文本消息
func replyTextMessage(message WechatmpMessage) (reply string, err error) {
	// 处理重试消息
	if message.Content == "重试" {
		reply = "未获取到历史消息，请尝试重新发送问题吧"
		if oldMsgId := MsgContext.Get(message.FromUserName).(int64); oldMsgId != 0 {
			out := MsgContext.Get(oldMsgId).(string)
			if out != "" {
				reply = out
				ezap.Debugf("获取到重试回复的消息: msgId: %d", oldMsgId)
			}
			ezap.Debugf("未获取到重试回复的消息: msgId: %d", oldMsgId)
		}
		return
	}

	// 处理 msgId 相同的重复请求，返回历史消息
	if MsgContext.Get(message.FromUserName).(int64) == message.MsgId {
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

		return
	}

	// 正常流程，询问 gemini
	reply, err = askGemini(message.Content)
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
