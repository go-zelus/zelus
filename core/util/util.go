package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// IsHTTPS 判断是否https协议
func IsHTTPS(request *http.Request) bool {
	schema := request.Header.Get("X-Forwarded-Proto")
	if schema == "https" || request.TLS != nil {
		return true
	}
	return false
}

// IsMobile 判断是否移动端
func IsMobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}
	isMobile := false
	mobileKeys := []string{"Mobile", "Android", "Silk/", "Kindle", "BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeys); i++ {
		if strings.Contains(userAgent, mobileKeys[i]) {
			isMobile = true
			break
		}
	}
	return isMobile
}

// ClientIP 获取客户端IP
func ClientIP(r *http.Request, c *websocket.Conn) string {
	ip := ""
	if r != nil {
		ip = r.Header.Get("X-Forwarded-For")
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		if ip == "" {
			ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		}
		if ip == "" && r.Host != "" {
			idx := strings.Index(r.Host, ":")
			ip = r.Host[:idx]
			if ip == "localhost" {
				ip = "127.0.0.1"
			}
		}
		if ip != "" {
			return ip
		}
	}
	if c != nil {
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.RemoteAddr().String())); err == nil {
			return ip
		}
		if ip == "" {
			return c.RemoteAddr().String()
		}
	}
	return ip
}

// MixMobile 混淆手机号
func MixMobile(mobile string) string {
	var phone string
	chars := strings.Split(mobile, "")
	if len(chars) < 7 {
		return mobile
	}
	for i := 0; i < len(chars); i++ {
		if i > 2 && i < 7 {
			phone += "*"
		} else {
			phone += chars[i]
		}
	}
	return phone
}

// MD5编码 32位小写
func MD5(message string) string {
	h := md5.New()
	h.Write([]byte(message))
	cipher := h.Sum(nil)
	return fmt.Sprintf("%x", cipher)
}

// MD5编码 16位小写
func MD5B16(message string) string {
	return MD5(message)[8:24]
}

// Base64 编码
func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Timestamp 毫秒时间戳
func Timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
