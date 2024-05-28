package wechatmp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fimreal/goutils/ezap"
)

// 公众号发送普通文本消息(未测试)
func SendWechatmpMessageToUser(message WechatmpMessage) error {
	// 获取 access_token
	accessToken, err := GetAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", accessToken)

	// 构造回复消息
	replyType := "text"
	replyContent := "当前公众号在维护中，请稍后再试。感谢您的关注！"
	reply := WechatmpMessage{
		ToUserName:   message.FromUserName,
		FromUserName: message.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      replyType,
		Content:      replyContent,
		MsgId:        message.MsgId,
	}
	ezap.Debugf("构造文本消息: %+v", reply)

	// 发送回复消息
	requestBody := new(bytes.Buffer)
	err = xml.NewEncoder(requestBody).Encode(reply)
	if err != nil {
		return err
	}
	client, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return err
	}
	client.Header.Set("Content-Type", "application/xml")

	response, err := http.DefaultClient.Do(client)
	if err != nil {
		return err
	}

	// 解析回复消息结果
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	ezap.Debugf("回复消息结果: %s, %s", response.Status, string(resBody))
	return nil
}
