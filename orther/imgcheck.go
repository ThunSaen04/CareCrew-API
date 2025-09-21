package orther

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

func IsImageFile(fh *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	default:
		return false
	}
}
