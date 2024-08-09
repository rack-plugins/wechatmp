package wechatmp

import (
	"net/http"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 接收公众号消息
// https://mp.weixin.qq.com/debug/cgi-bin/apiinfo?t=index&type=%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8F%9C%E5%8D%95&form=%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8F%9C%E5%8D%95%E5%88%9B%E5%BB%BA%E6%8E%A5%E5%8F%A3%20/menu/creat
func HandleRequest(c *gin.Context) {
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
	ezap.Debugf("收到请求信息: %+v", message)

	// 生成回复消息
	reply, err := PassiveReply(message)
	if err != nil {
		ezap.Error("构造回复消息时出错: ", err.Error())
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ezap.Debugf("回复消息: %+v", reply)

	c.XML(http.StatusOK, reply)
}

func PassiveReply(m WechatmpMessage) (*WechatmpMessage, error) {
	// 根据消息类型回复
	switch m.MsgType {
	case "text":
		return m.HandleText()
	case "event":
		return m.HandleEvent()
	// case "image":
	default:
		return &WechatmpMessage{
			ToUserName:   m.FromUserName,
			FromUserName: m.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgId:        m.MsgId,
			MsgType:      "text",
			Content:      "我现在无法处理此类型消息",
		}, nil
	}

}
