package utils

import (
	"fmt"
	"strings"
)

// GetAvatarURL 获取头像完整URL
func GetAvatarURL(avatar string) string {
	// 如果头像为空，返回默认头像
	if avatar == "" {
		return "/static/avatar/default.png"
	}

	// 如果已经是完整URL（包含http或https），直接返回
	if strings.HasPrefix(avatar, "http://") || strings.HasPrefix(avatar, "https://") {
		return avatar
	}

	// 如果是相对路径，拼接静态资源地址
	return fmt.Sprintf("/static/avatar/%s", avatar)
}

// GetAvatarFullURL 获取头像完整URL（包含域名）
func GetAvatarFullURL(avatar string, baseURL string) string {
	relativeURL := GetAvatarURL(avatar)

	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	// 拼接完整URL
	return fmt.Sprintf("%s%s", baseURL, relativeURL)
}
