// 微信素材库接口
// 用于上传图片、视频、语音、缩略图等素材，在消息中引用
package wechatmp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 微信公众号素材库上传地址，分为临时3天存储和永久存储
// 图片（image）: 10M，支持bmp/png/jpeg/jpg/gif格式
// 语音（voice）：2M，播放长度不超过60s，mp3/wma/wav/amr格式
// 视频（video）：10MB，支持MP4格式
// 缩略图（thumb）：64KB，支持JPG格式
// uploadImage 接口所上传的图片不占用公众号的素材库中图片数量的100000个的限制。图片仅支持jpg/png格式，大小必须在1MB以下。这里未实现
const (
	uploadTempURL      = "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	uploadPermanentURL = "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%s&type=%s"
	// uploadImageURL     = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s"
)

// VideoDescription 视频描述结构体
type VideoDescription struct {
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
}

// MediaResponse 上传媒体文件后的响应结构
type MediaResponse struct {
	Type      string `json:"type,omitempty"`       // 媒体文件类型
	MediaID   string `json:"media_id,omitempty"`   // 媒体文件标识
	CreatedAt int64  `json:"created_at,omitempty"` // 媒体文件上传时间戳
	// URL       string `json:"url,omitempty"`        // 图片素材的 URL, 仅针对图片 uploadImage 接口，这里未实现
}

// DownloadFile 下载文件并保存到指定目录，返回文件路径
func DownloadFile(url, downloadDir string) (string, error) {
	// 设置HTTP客户端，添加超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发起GET请求
	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求文件[%s]失败: %w", url, err)
	}
	defer response.Body.Close()

	// 检查响应状态
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取文件[%s]失败，状态码: %d", url, response.StatusCode)
	}

	// 安全地获取文件名并构建完整路径
	fileName := filepath.Base(url)
	filePath := filepath.Join(downloadDir, fileName)

	// 检查并创建下载目录
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建目录[%s]失败: %w", downloadDir, err)
	}

	// 创建文件并处理潜在错误
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件[%s]失败: %w", filePath, err)
	}
	defer file.Close()

	// 将响应内容写入文件
	if _, err := io.Copy(file, response.Body); err != nil {
		return "", fmt.Errorf("写入文件[%s]失败: %w", filePath, err)
	}

	return filePath, nil
}

// UploadTempMedia 上传临时多媒体文件到微信公共平台
func UploadTempMedia(mediaType, filePath string) (*MediaResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 创建一个新的表单文件
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)

	// 添加媒体文件到表单
	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建form文件失败: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 关闭写入器以完成表单
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭writer失败: %v", err)
	}

	// 构建请求
	accessToken := WechatmpAccessToken.Token
	if accessToken == "" {
		return nil, fmt.Errorf("获取access token失败: %v", err)
	}
	request, err := http.NewRequest("POST", fmt.Sprintf(uploadTempURL, accessToken, mediaType), bodyBuf)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer response.Body.Close()

	// 读取响应
	var mediaResp MediaResponse
	if err := json.NewDecoder(response.Body).Decode(&mediaResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &mediaResp, nil
}

// UploadPermanentMedia 上传永久多媒体文件到微信公共平台
func UploadPermanentMedia(mediaType, filePath string, videoDesc *VideoDescription) (*MediaResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 创建一个新的表单文件
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)

	// 添加媒体文件到表单
	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建form文件失败: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 处理视频描述，如果有的话
	if mediaType == "video" && videoDesc != nil {
		descPart, err := writer.CreateFormField("description")
		if err != nil {
			return nil, fmt.Errorf("创建描述字段失败: %v", err)
		}
		descJSON, err := json.Marshal(videoDesc)
		if err != nil {
			return nil, fmt.Errorf("序列化视频描述失败: %v", err)
		}
		_, err = descPart.Write(descJSON)
		if err != nil {
			return nil, fmt.Errorf("写入描述字段失败: %v", err)
		}
	}

	// 关闭写入器以完成表单
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭writer失败: %v", err)
	}

	// 构建请求
	accessToken := WechatmpAccessToken.Token
	if accessToken == "" {
		return nil, fmt.Errorf("获取access token失败: %v", err)
	}
	request, err := http.NewRequest("POST", fmt.Sprintf(uploadPermanentURL, accessToken, mediaType), bodyBuf)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer response.Body.Close()

	// 读取响应
	var mediaResp MediaResponse
	if err := json.NewDecoder(response.Body).Decode(&mediaResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &mediaResp, nil
}
