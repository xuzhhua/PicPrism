package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaevor/go-nanoid"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

type ProcessResult struct {
	ID       string
	Ext      string
	Size     int64
	Width    int
	Height   int
	MimeType string
	Hash     string
}

type Links struct {
	URL      string `json:"url"`
	WebPURL  string `json:"webp_url"`
	ThumbURL string `json:"thumb_url"`
	Markdown string `json:"markdown"`
	HTML     string `json:"html"`
	BBCode   string `json:"bbcode"`
}

var newID func() string

func init() {
	gen, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}
	newID = gen
}

// ProcessUpload 处理上传的图片文件，返回处理结果
func ProcessUpload(r io.Reader, originalName, dataDir string) (*ProcessResult, error) {
	// 读取全部内容
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// 计算 SHA-256
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])

	// 检测格式
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext == "" {
		ext = ".bin"
	}
	mimeType, ext := detectFormat(data, ext)

	id := newID()

	// 存原图
	origPath := filepath.Join(dataDir, "images", id+ext)
	if err := os.WriteFile(origPath, data, 0644); err != nil {
		return nil, fmt.Errorf("write original: %w", err)
	}

	// 解析图片尺寸（SVG 跳过）
	width, height := 0, 0
	if ext != ".svg" {
		img, _, decErr := decodeImage(data, ext)
		if decErr == nil && img != nil {
			bounds := img.Bounds()
			width = bounds.Dx()
			height = bounds.Dy()

			// 同步生成优化版和缩略图（JPEG 格式），确保响应返回前文件已就绪
			generateOptimized(img, filepath.Join(dataDir, "images", id+"_opt.jpg"))
			generateThumb(img, filepath.Join(dataDir, "thumbs", id+"_thumb.jpg"), 800)
		}
	}

	return &ProcessResult{
		ID:       id,
		Ext:      ext,
		Size:     int64(len(data)),
		Width:    width,
		Height:   height,
		MimeType: mimeType,
		Hash:     hash,
	}, nil
}

// BuildLinks 根据图片 ID、扩展名和 base URL 生成各格式链接
func BuildLinks(id, ext, baseURL string) Links {
	origURL := fmt.Sprintf("%s/i/%s%s", baseURL, id, ext)
	webpURL := fmt.Sprintf("%s/i/%s_opt.jpg", baseURL, id)
	thumbURL := fmt.Sprintf("%s/i/%s_thumb.jpg", baseURL, id)

	// SVG 不转换，直链回退
	if ext == ".svg" {
		webpURL = origURL
		thumbURL = origURL
	}

	return Links{
		URL:      origURL,
		WebPURL:  webpURL,
		ThumbURL: thumbURL,
		Markdown: fmt.Sprintf("![image](%s)", webpURL),
		HTML:     fmt.Sprintf(`<img src="%s" alt="image">`, webpURL),
		BBCode:   fmt.Sprintf("[img]%s[/img]", webpURL),
	}
}

// detectFormat 根据文件头和扩展名返回 MIME 类型和规范扩展名
func detectFormat(data []byte, ext string) (string, string) {
	if len(data) >= 4 {
		switch {
		case data[0] == 0xFF && data[1] == 0xD8:
			return "image/jpeg", ".jpg"
		case data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G':
			return "image/png", ".png"
		case data[0] == 'G' && data[1] == 'I' && data[2] == 'F':
			return "image/gif", ".gif"
		case len(data) >= 12 && string(data[8:12]) == "WEBP":
			return "image/webp", ".webp"
		case data[0] == 'B' && data[1] == 'M':
			return "image/bmp", ".bmp"
		case (data[0] == 'I' && data[1] == 'I') || (data[0] == 'M' && data[1] == 'M'):
			return "image/tiff", ".tiff"
		case len(data) >= 12 && (string(data[4:8]) == "ftyp"):
			return "image/avif", ".avif"
		}
	}
	// SVG 检测
	s := strings.TrimSpace(string(data[:min(100, len(data))]))
	if strings.HasPrefix(s, "<svg") || strings.HasPrefix(s, "<?xml") || ext == ".svg" {
		return "image/svg+xml", ".svg"
	}
	return "application/octet-stream", ext
}

func decodeImage(data []byte, ext string) (image.Image, string, error) {
	reader := bytes.NewReader(data)

	switch ext {
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(reader)
		return img, "jpeg", err
	case ".png":
		img, err := png.Decode(reader)
		return img, "png", err
	case ".gif":
		img, err := gif.Decode(reader)
		return img, "gif", err
	case ".webp":
		img, err := webp.Decode(reader)
		return img, "webp", err
	case ".bmp":
		img, err := bmp.Decode(reader)
		return img, "bmp", err
	case ".tiff", ".tif":
		img, err := tiff.Decode(reader)
		return img, "tiff", err
	default:
		// 尝试标准库自动检测
		img, format, err := image.Decode(reader)
		return img, format, err
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
