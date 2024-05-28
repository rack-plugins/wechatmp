package wechatmp

import (
	"encoding/xml"
	"sync"
	"time"

	"github.com/gogf/gf/container/gmap"
)

var (
	// 存储消息
	MsgContext = gmap.New()
	// 存储 access_token，用于请求订阅号 api
	WechatmpAccessToken AccessToken
)

type AccessToken struct {
	Token  string
	Mutex  sync.RWMutex
	Expiry time.Time
}

// signature	微信加密签名，signature结合了开发者填写的token参数和请求中的timestamp参数、nonce参数。
// timestamp	时间戳
// nonce	随机数
// echostr	随机字符串
type Signature struct {
	Signature     string `json:"signature" form:"signature" xml:"signature" validate:"required"`
	Timestamp     string `json:"timestamp" form:"timestamp" xml:"timestamp" validate:"required"`
	Nonce         string `json:"nonce" form:"nonce" xml:"nonce" validate:"required"`
	Echostr       string `json:"echostr" form:"echostr" xml:"echostr" `
	Openid        string `json:"openid" form:"openid" xml:"openid" `
	Encrtpt_type  string `json:"encrtpt_type" form:"encrtpt_type" xml:"encrtpt_type" `
	Msg_signature string `json:"msg_signature" form:"msg_signature" xml:"msg_signature" `
}

// type Cdata struct {
// 	Value string `xml:",cdata"`
// }

// 微信公众号消息结构体
type WechatmpMessage struct {
	XMLName xml.Name `json:"-" xml:"xml"` // 指定 xml 根标签

	ToUserName   string `json:"ToUserName" xml:"ToUserName"`
	FromUserName string `json:"FromUserName" xml:"FromUserName"`
	CreateTime   int64  `json:"CreateTime" xml:"CreateTime"`
	MsgType      string `json:"MsgType" xml:"MsgType"` // text, image, voice, video, shortvideo, location, link

	// text
	Content string `json:"Content" xml:"Content"`
	MsgId   int64  `json:"MsgId" xml:"MsgId"`

	MediaId string `json:"MediaId" xml:"MediaId"` // image, voice, video, shortvideo
	// image
	PicUrl string `json:"PicUrl" xml:"PicUrl"`
	// voice
	Format string `json:"Format" xml:"Format"` // amr, speex
	// video shortvideo
	ThumbMediaId string `json:"ThumbMediaId" xml:"ThumbMediaId"`

	// location
	Location_X float64 `json:"Location_X" xml:"Location_X"`
	Location_Y float64 `json:"Location_Y" xml:"Location_Y"`
	Scale      int64   `json:"Scale" xml:"Scale"`
	Label      string  `json:"Label" xml:"Label"`
	// link
	Title       string `json:"Title" xml:"Title"`
	Description string `json:"Description" xml:"Description"`
	URL         string `json:"URL" xml:"URL"`
}
