package chat

import (
	"fmt"
	"gochat/internal/pkg/code_msg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传文件（图片/文件）
func (h *ChatHandler) UploadFile(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		h.response.JsonError(ctx, err, "获取上传文件失败")
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

	// 返回文件信息
	h.response.JsonSuccess(ctx, gin.H{
		"file_url":  fileInfo.URL,
		"file_name": fileInfo.Name,
		"file_size": fileInfo.Size,
		"file_type": fileInfo.Type,
	})
}

// FileInfo 文件信息
type FileInfo struct {
	URL  string `json:"file_url"`
	Name string `json:"file_name"`
	Size int64  `json:"file_size"`
	Type string `json:"file_type"`
}

// validateFile 验证文件
func (h *ChatHandler) validateFile(header *multipart.FileHeader) (string, code_msg.BusinessCode, error) {
	// 检查文件大小（限制为10MB）
	const maxSize = 10 * 1024 * 1024
	if header.Size > maxSize {
		return "", code_msg.BadRequest, fmt.Errorf("文件大小不能超过10MB")
	}

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	ext := filepath.Ext(header.Filename)

	switch {
	case isImageFile(contentType, ext):
		return "image", 0, nil
	case isDocumentFile(contentType, ext):
		return "document", 0, nil
	default:
		return "", code_msg.BadRequest, fmt.Errorf("不支持的文件类型")
	}
}

// saveFile 保存文件
func (h *ChatHandler) saveFile(file multipart.File, header *multipart.FileHeader, fileType string) (*FileInfo, code_msg.BusinessCode, error) {
	// 生成文件名
	fileName := generateFileName(header.Filename)

	// 创建目录
	uploadDir := fmt.Sprintf("uploads/%s/%s", fileType, time.Now().Format("2006/01/02"))
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

	// 返回文件信息
	fileInfo := &FileInfo{
		URL:  "/" + filePath, // 相对路径，实际部署时需要配置静态文件服务
		Name: header.Filename,
		Size: header.Size,
		Type: fileType,
	}

	return fileInfo, 0, nil
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

// generateFileName 生成文件名
func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s%s", timestamp, strconv.FormatInt(time.Now().Unix(), 36), ext)
}
