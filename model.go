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

// Signature represents the structure for WeChat message verification.
type Signature struct {
	Signature    string `json:"signature" form:"signature" xml:"signature" validate:"required"` // WeChat encrypted signature
	Timestamp    string `json:"timestamp" form:"timestamp" xml:"timestamp" validate:"required"` // Timestamp of the request
	Nonce        string `json:"nonce" form:"nonce" xml:"nonce" validate:"required"`             // Random number for validation
	Echostr      string `json:"echostr" form:"echostr" xml:"echostr"`                           // Random string used during verification
	Openid       string `json:"openid" form:"openid" xml:"openid"`                              // User's unique identifier
	EncryptType  string `json:"encrypt_type" form:"encrypt_type" xml:"encrypt_type"`            // Type of encryption used (e.g., aes)
	MsgSignature string `json:"msg_signature" form:"msg_signature" xml:"msg_signature"`         // Message signature for verifying the integrity
}

// WechatmpMessage represents the structure of a WeChat public account message.
// 个人公众号不支持开通客服功能，仅使用被动回复功能不多
type WechatmpMessage struct {
	XMLName xml.Name `json:"-" xml:"xml"` // Specify the XML root tag

	ToUserName   string `json:"ToUserName" xml:"ToUserName"`     // Recipient's user ID
	FromUserName string `json:"FromUserName" xml:"FromUserName"` // Sender's public account ID
	CreateTime   int64  `json:"CreateTime" xml:"CreateTime"`     // Message creation time (timestamp)
	MsgType      string `json:"MsgType" xml:"MsgType"`           // Type of message: event, text, image, voice, video, shortvideo, location, link
	MsgId        int64  `json:"MsgId" xml:"MsgId"`               // Message ID

	// Event messages
	Event    string `json:"Event,omitempty" xml:"Event,omitempty"`       // Event type: subscribe, unsubscribe, SCAN, LOCATION, CLICK, VIEW
	EventKey string `json:"EventKey,omitempty" xml:"EventKey,omitempty"` // Key value associated with the event, such as QR code parameters or custom menu keys

	// Text messages
	Content string `json:"Content,omitempty" xml:"Content,omitempty"` // Content of the text message

	// Media messages
	Image *Image `json:"Image,omitempty" xml:"Image,omitempty"` // Image media information
	Voice *Voice `json:"Voice,omitempty" xml:"Voice,omitempty"` // Voice media information
	Video *Video `json:"Video,omitempty" xml:"Video,omitempty"` // Video media information
	Music *Music `json:"Music,omitempty" xml:"Music,omitempty"` // Music media information

	// News messages
	// ArticleCount int       `json:"ArticleCount,omitempty" xml:"ArticleCount,omitempty"` // Number of articles in the news message
	// Articles     *Articles `json:"Articles,omitempty" xml:"Articles,omitempty"`         // Articles contained in the news message
}

// Image represents an image media structure.
type Image struct {
	MediaId string `json:"MediaId" xml:"MediaId" cdata:",chardata"` // Media ID of the image
}

// Voice represents a voice media structure.
type Voice struct {
	MediaId string `json:"MediaId" xml:"MediaId" cdata:",chardata"` // Media ID of the voice
}

// Video represents a video media structure.
type Video struct {
	MediaId     string `json:"MediaId" xml:"MediaId" cdata:",chardata"`         // Media ID of the video
	Title       string `json:"Title" xml:"Title" cdata:",chardata"`             // Title of the video
	Description string `json:"Description" xml:"Description" cdata:",chardata"` // Description of the video
}

// Music represents a music media structure.
type Music struct {
	Title        string `json:"Title" xml:"Title" cdata:",chardata"`               // Title of the music
	Description  string `json:"Description" xml:"Description" cdata:",chardata"`   // Description of the music
	MusicUrl     string `json:"MusicUrl" xml:"MusicUrl" cdata:",chardata"`         // URL of the music
	HQMusicUrl   string `json:"HQMusicUrl" xml:"HQMusicUrl" cdata:",chardata"`     // High-quality music URL
	ThumbMediaId string `json:"ThumbMediaId" xml:"ThumbMediaId" cdata:",chardata"` // Media ID of the thumbnail
}

// Articles represents a collection of articles in a news message.
type Articles struct {
	Items []Item `json:"Items" xml:"Articles>item"` // List of articles
}

// Item represents an individual article in a news message.
type Item struct {
	Title       string `json:"Title" xml:"Title" cdata:",chardata"`             // Title of the article
	Description string `json:"Description" xml:"Description" cdata:",chardata"` // Description of the article
	PicUrl      string `json:"PicUrl" xml:"PicUrl" cdata:",chardata"`           // Picture URL of the article
	Url         string `json:"Url" xml:"Url" cdata:",chardata"`                 // URL of the article
}
