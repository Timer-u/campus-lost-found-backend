package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"campus-lost-found-backend/pkg/response"
)

const (
	uploadDir     = "uploads"
	maxImageBytes = 5 << 20 // 5MB，文档定义的单张图片上限
)

// SaveImage 校验并保存上传的图片，返回可公开访问的相对 URL（/uploads/xxx）
func SaveImage(fh *multipart.FileHeader) (string, *response.Errno) {
	if fh.Size > maxImageBytes {
		return "", response.ErrImageTooLarge
	}

	src, err := fh.Open()
	if err != nil {
		return "", response.ErrInternal
	}
	defer src.Close()

	// 内容嗅探：只接受真实图片，防止改扩展名绕过
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return "", response.ErrInternal
	}
	if ct := http.DetectContentType(buf[:n]); !strings.HasPrefix(ct, "image/") {
		return "", response.ErrInvalidParams.WithMsg("仅支持上传图片文件")
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", response.ErrInvalidParams.WithMsg("仅支持 jpg/jpeg/png/gif/webp 图片")
	}

	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		return "", response.ErrInternal
	}
	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(random), ext)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", response.ErrInternal
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", response.ErrInternal
	}
	dst, err := os.Create(filepath.Join(uploadDir, name))
	if err != nil {
		return "", response.ErrInternal
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", response.ErrInternal
	}
	return "/uploads/" + name, nil
}
