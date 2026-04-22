package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	dataDir string
}

func NewFileHandler(dataDir string) *FileHandler {
	return &FileHandler{dataDir: dataDir}
}

// Serve GET /i/:filename
// 支持：<id>.<ext>  |  <id>.webp  |  <id>_thumb.webp
func (h *FileHandler) Serve(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" || strings.ContainsAny(filename, "/\\..") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid filename"})
	}

	// 安全检查：防止路径穿越
	var diskPath string
	if strings.HasSuffix(filename, "_thumb.webp") {
		diskPath = filepath.Join(h.dataDir, "thumbs", filename)
	} else {
		diskPath = filepath.Join(h.dataDir, "images", filename)
	}

	// 确保路径在数据目录内
	absData, _ := filepath.Abs(h.dataDir)
	absPath, _ := filepath.Abs(diskPath)
	if !strings.HasPrefix(absPath, absData) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	// 设置缓存头
	c.Set("Cache-Control", "public, max-age=31536000, immutable")
	return c.SendFile(diskPath)
}
