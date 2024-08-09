package wechatmp

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CheckSignature(c *gin.Context) {
	var sig Signature
	if err := c.ShouldBind(&sig); err != nil {
		ezap.Error(err.Error())
		c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sig.Signature == "" {
		sig.Echostr = c.Query("echostr")
		sig.Nonce = c.Query("nonce")
		sig.Signature = c.Query("signature")
		sig.Timestamp = c.Query("timestamp")
	}

	token := viper.GetString(ID + ".token")
	if sig.checkSignature(token) {
		c.String(http.StatusOK, sig.Echostr)
	} else {
		c.String(http.StatusForbidden, "signature error.")
	}
}

// checkSignature verifies the signature against the given timestamp and nonce.
// 参考：https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
// 调用测试：https://developers.weixin.qq.com/apiExplorer?type=messagePush
func (sig Signature) checkSignature(token string) bool {
	tmpArr := []string{token, sig.Timestamp, sig.Nonce}
	sort.Strings(tmpArr)

	tmpStr := strings.Join(tmpArr, "")

	tmpStrHash := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))
	ezap.Debugf("challenge signature: %s, answer: %s", sig.Signature, tmpStrHash)

	return tmpStrHash == sig.Signature
}

// GetAccessToken 获取 access_token
func GetAccessToken() (accessToken string, err error) {
	// 在并发情况下，多个 goroutine 可能同时访问 WechatmpAccessToken，读锁（RLock）允许多个 goroutine 同时读取数据，但阻止任何 goroutine 获取写锁。
	WechatmpAccessToken.Mutex.RLock()
	defer WechatmpAccessToken.Mutex.RUnlock()

	// 检查 access_token 是否过期
	if WechatmpAccessToken.Expiry.Before(time.Now()) {
		// 如果过期，则刷新 access_token
		err := refreshAccessToken()
		if err != nil {
			return "", err
		}
	}

	return WechatmpAccessToken.Token, nil
}

func refreshAccessToken() error {
	appid, appsecret := viper.GetString(ID+".appid"), viper.GetString(ID+".appsecret")

	// 获取 access_token
	accessToken, err := requestAccessToken(appid, appsecret)
	if err != nil {
		return fmt.Errorf("刷新公众号 access_token 失败: %v", err)
	}

	// 打印 access_token，方便 debug
	ezap.Debugf("刷新公众号 access_token: %s", accessToken)

	// 更新全局缓存
	WechatmpAccessToken.Mutex.Lock()
	WechatmpAccessToken.Token = accessToken
	WechatmpAccessToken.Expiry = time.Now().Add(time.Hour * 1) // 设置过期时间为1小时, access_token 有效期2小时
	WechatmpAccessToken.Mutex.Unlock()

	ezap.Infof("公众号 access_token 已刷新, 过期时间: %v", WechatmpAccessToken.Expiry)

	return nil
}

// 参考：https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func requestAccessToken(appid, appsecret string) (accessToken string, err error) {
	const clientCredentialURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	url := fmt.Sprintf(clientCredentialURL, appid, appsecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if res.Errcode != 0 {
		return "", fmt.Errorf("WeChat API error: code=%d, msg=%s", res.Errcode, res.Errmsg)
	}

	return res.AccessToken, nil
}
