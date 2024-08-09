package wechatmp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fimreal/goutils/ezap"
)

const helpMsg = `
这是帮助信息，目前支持的指令有：
/help 帮助，显示这条信息，展示我可以接受的指令
/ask 提问，从数据库中查询并支持回复你的问题。（默认对话能力）
/draw 绘画，后面添加文字描述，我会为你绘制一张符合的图片`

// 被动回复消息 - 处理文本消息
func (m *WechatmpMessage) HandleText() (reply *WechatmpMessage, err error) {
	// 创建回复消息模板
	reply = &WechatmpMessage{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgId:        m.MsgId,
		MsgType:      "text", // default: text
	}

	// 处理重试消息
	if m.Content == "重试" {
		reply.Content = "未获取到历史消息，请尝试重新发送问题吧"
		if oldMsgId := MsgContext.GetOrSet(m.FromUserName, int64(0)).(int64); oldMsgId != int64(0) {
			out := MsgContext.Get(oldMsgId).(*WechatmpMessage)
			if out != nil {
				reply = out
				ezap.Infof("获取到重试回复的消息: msgId: %d", oldMsgId)
			} else {
				ezap.Infof("未获取到重试回复的消息: msgId: %d", oldMsgId)
			}
		}
		return
	}

	// 处理 msgId 相同的重复请求，返回历史消息
	// 微信服务器在将用户的消息发给公众号的开发者服务器地址（开发者中心处配置）后，微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次。
	if MsgContext.GetOrSet(m.FromUserName, int64(0)).(int64) == m.MsgId {
		ezap.Infof("获取到超时回复的消息: msgId: %d", m.MsgId)

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		replyChan := make(chan *WechatmpMessage, 1)
		go func() {
			for {
				reply := MsgContext.Get(m.MsgId)
				if reply != nil {
					replyChan <- reply.(*WechatmpMessage)
					return
				}
				time.Sleep(100 * time.Microsecond)
			}
		}()

		select {
		case reply = <-replyChan:
			return
		case <-ctx.Done():
			reply.Content = "生成内容需要一点时间，请稍后发送【重试】重新获取结果。"
			return
		}

	}

	// 保存新消息到用户对话上下文
	MsgContext.Set(m.FromUserName, m.MsgId)

	// 识别指令
	switch {
	case strings.HasPrefix(m.Content, "/help"):
		reply.Content = helpMsg
		return
	case strings.HasPrefix(m.Content, "/draw"):
		m.Content = strings.TrimPrefix(m.Content, "/draw")
		// 绘画
		media, e := m.txt2ImgReply()
		if e != nil {
			ezap.Errorf("%w", e)
			reply.Content = "我好像处理不过来了，请稍后再试再试一下吧"
		} else {
			// 发送图片消息
			reply.MsgType = "image"
			reply.Image = new(Image)
			reply.Image.MediaId = media.MediaID
		}

	// case strings.HasPrefix(m.Content, "/ask"):
	default:
		// 过滤 /ask 指令
		m.Content = strings.TrimPrefix(m.Content, "/ask")
		reply.Content, err = m.llmReply()
		if err != nil {
			ezap.Errorf("llm ask [%s] error: %v", m.Content, err)
			reply.Content = "抱歉，我出了点问题，请稍后再试或者换一个方式提问"
		}
	}
	// 保存新消息到用户对话上下文
	MsgContext.Set(m.MsgId, reply)
	// 15 分钟后删除历史消息，减少内存占用
	go func() {
		time.Sleep(15 * time.Minute)
		MsgContext.Remove(m.MsgId)
	}()
	return
}

// 构造回复消息 - 大模型回复
func (m *WechatmpMessage) llmReply() (reply string, err error) {
	// 正常流程，询问 llm
	return LLM.Ask(m.Content)
}

// 构造回复消息 - 绘画
func (m *WechatmpMessage) txt2ImgReply() (media *MediaResponse, err error) {
	picUrl, e := Txt2Anime(m.Content)
	if e != nil {
		err = fmt.Errorf("mlm draw [%s] error: %v", m.Content, e)
		return
	}
	// 下载图片
	picPath, e := DownloadFile(picUrl, ".cache")
	if e != nil {
		err = fmt.Errorf("download picurl[%s] error: %v", picUrl, e)
		return
	}
	// 上传存储为临时素材
	return UploadTempMedia("image", picPath)
}
