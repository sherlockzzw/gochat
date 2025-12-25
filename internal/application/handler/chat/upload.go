package chat

import (
	"fmt"
	"gochat/api/api/chat"
	"gochat/internal/pkg/code_msg"
	globalUtils "gochat/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传文件（支持图片/文档/视频/语音/表情包）
func (h *ChatHandler) UploadFile(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		h.response.JsonErrorFixation(ctx, code_msg.BadRequest)
		return
	}
	defer file.Close()

	// 验证文件类型和大小
	fileType, code, err := h.validateFile(header)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	// 保存文件
	fileInfo, code, err := h.saveFile(file, header, fileType)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	// 获取完整URL
	apiPort := 8080
	if port := globalUtils.GetApiPort(); port > 0 {
		apiPort = port
	}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", apiPort)
	fullURL := h.getFileFullURL(fileInfo.RelativePath, baseURL)

	// 构造protobuf响应
	resp := &chat.UploadFileResponse{
		FileUrl:  fullURL,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
	}

	// 返回文件信息
	h.response.JsonSuccess(ctx, resp)
}

// FileInfo 文件信息
type FileInfo struct {
	RelativePath string // 相对路径（用于构造完整URL）
	Name         string // 文件名
	Size         int64  // 文件大小
	Type         string // 文件类型
}

// validateFile 验证文件
func (h *ChatHandler) validateFile(header *multipart.FileHeader) (string, code_msg.BusinessCode, error) {
	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	ext := strings.ToLower(filepath.Ext(header.Filename))

	var maxSize int64
	var fileType string

	switch {
	case isImageFile(contentType, ext):
		maxSize = 10 * 1024 * 1024 // 10MB
		fileType = "image"
	case isDocumentFile(contentType, ext):
		maxSize = 10 * 1024 * 1024 // 10MB
		fileType = "document"
	case isVideoFile(contentType, ext):
		maxSize = 100 * 1024 * 1024 // 100MB（视频文件较大）
		fileType = "video"
	case isVoiceFile(contentType, ext):
		maxSize = 10 * 1024 * 1024 // 10MB（语音文件）
		fileType = "voice"
	case isEmojiFile(contentType, ext):
		maxSize = 5 * 1024 * 1024 // 5MB（表情包文件）
		fileType = "emoji"
	default:
		return "", code_msg.BadRequest, fmt.Errorf("不支持的文件类型")
	}

	// 检查文件大小
	if header.Size > maxSize {
		return "", code_msg.BadRequest, fmt.Errorf("文件大小不能超过%dMB", maxSize/(1024*1024))
	}

	return fileType, 0, nil
}

// saveFile 保存文件
func (h *ChatHandler) saveFile(file multipart.File, header *multipart.FileHeader, fileType string) (*FileInfo, code_msg.BusinessCode, error) {
	// 生成文件名
	fileName := generateFileName(header.Filename)

	// 创建目录（存储到resources目录下）
	uploadDir := fmt.Sprintf("resources/upload/%s/%s", fileType, time.Now().Format("2006/01/02"))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, code_msg.ServerError, fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 文件路径
	filePath := filepath.Join(uploadDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, code_msg.ServerError, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, code_msg.ServerError, fmt.Errorf("保存文件失败: %w", err)
	}

	// 构造相对路径（用于静态资源访问）
	relativePath := fmt.Sprintf("/static/upload/%s/%s/%s", fileType, time.Now().Format("2006/01/02"), fileName)

	// 返回文件信息
	fileInfo := &FileInfo{
		RelativePath: relativePath,
		Name:         header.Filename,
		Size:         header.Size,
		Type:         fileType,
	}

	return fileInfo, 0, nil
}

// getFileFullURL 获取文件完整URL
func (h *ChatHandler) getFileFullURL(relativePath string, baseURL string) string {
	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	// 拼接完整URL
	return fmt.Sprintf("%s%s", baseURL, relativePath)
}

// isImageFile 检查是否为图片文件
func isImageFile(contentType, ext string) bool {
	imageTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	return imageTypes[contentType] || imageExts[ext]
}

// isDocumentFile 检查是否为文档文件
func isDocumentFile(contentType, ext string) bool {
	docTypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"text/plain": true,
	}

	docExts := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
	}

	return docTypes[contentType] || docExts[ext]
}

// isVideoFile 检查是否为视频文件
func isVideoFile(contentType, ext string) bool {
	videoTypes := map[string]bool{
		"video/mp4":             true,
		"video/mpeg":            true,
		"video/quicktime":        true,
		"video/x-msvideo":        true,
		"video/x-ms-wmv":         true,
		"video/webm":             true,
		"video/x-flv":            true,
		"video/3gpp":             true,
		"application/vnd.apple.mpegurl": true, // HLS
		"application/x-mpegURL": true,          // HLS
	}

	videoExts := map[string]bool{
		".mp4":  true,
		".mpeg": true,
		".mpg":  true,
		".mov":  true,
		".avi":  true,
		".wmv":  true,
		".flv":  true,
		".webm": true,
		".mkv":  true,
		".3gp":  true,
		".3gpp": true,
		".m3u8": true, // HLS播放列表
	}

	return videoTypes[contentType] || videoExts[ext]
}

// isVoiceFile 检查是否为语音文件
func isVoiceFile(contentType, ext string) bool {
	voiceTypes := map[string]bool{
		"audio/mpeg":      true,
		"audio/mp3":       true,
		"audio/wav":       true,
		"audio/wave":      true,
		"audio/x-wav":     true,
		"audio/ogg":       true,
		"audio/vorbis":    true,
		"audio/aac":       true,
		"audio/mp4":       true,
		"audio/x-m4a":     true,
		"audio/amr":       true,
		"audio/webm":      true,
		"application/octet-stream": true, // 某些音频文件可能使用此类型
	}

	voiceExts := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".ogg":  true,
		".oga":  true,
		".aac":  true,
		".m4a":  true,
		".amr":  true,
		".webm": true,
		".wma":  true,
		".flac": true,
	}

	return voiceTypes[contentType] || voiceExts[ext]
}

// isEmojiFile 检查是否为表情包文件
func isEmojiFile(contentType, ext string) bool {
	// 表情包可以是GIF动画、静态图片等
	emojiTypes := map[string]bool{
		"image/gif":      true,
		"image/jpeg":    true,
		"image/jpg":     true,
		"image/png":     true,
		"image/webp":    true,
		"image/apng":    true, // 动画PNG
	}

	emojiExts := map[string]bool{
		".gif":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".apng": true,
	}

	return emojiTypes[contentType] || emojiExts[ext]
}

// generateFileName 生成文件名
func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s%s", timestamp, strconv.FormatInt(time.Now().Unix(), 36), ext)
}
