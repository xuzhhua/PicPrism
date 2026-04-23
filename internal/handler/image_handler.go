package handler

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/picprism/picprism/internal/service"
	"github.com/picprism/picprism/internal/store"
)

type ImageHandler struct {
	images  *store.ImageStore
	tags    *store.TagStore
	dataDir string
}

func NewImageHandler(images *store.ImageStore, tags *store.TagStore, dataDir string) *ImageHandler {
	return &ImageHandler{images: images, tags: tags, dataDir: dataDir}
}

// Upload POST /api/v1/images
func (h *ImageHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file field required"})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot open file"})
	}
	defer f.Close()

	result, err := service.ProcessUpload(f, file.Filename, h.dataDir)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// 检查重复（hash 冲突）
	existing, err := h.images.GetByHash(result.Hash)
	if err == nil && existing != nil {
		// 已存在，直接返回现有记录的链接
		baseURL := service.ExtractBaseURL(c)
		tags, _ := h.tags.GetImageTags(existing.ID)
		tagNames := make([]string, len(tags))
		for i, t := range tags {
			tagNames[i] = t.Name
		}
		links := service.BuildLinks(existing.ID, existing.Ext, baseURL)
		return c.Status(fiber.StatusOK).JSON(buildImageResponse(existing, tagNames, links))
	}

	// 解析 tags 参数（逗号分隔）
	tagNames := parseTags(c.FormValue("tags"))
	var tagIDs []int
	if len(tagNames) > 0 {
		savedTags, err := h.tags.UpsertTags(tagNames)
		if err == nil {
			for _, t := range savedTags {
				tagIDs = append(tagIDs, t.ID)
			}
		}
	}

	img := &store.Image{
		ID:        result.ID,
		Filename:  file.Filename,
		Ext:       result.Ext,
		Size:      result.Size,
		Width:     result.Width,
		Height:    result.Height,
		MimeType:  result.MimeType,
		Hash:      result.Hash,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.images.Create(img); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db write failed"})
	}

	if len(tagIDs) > 0 {
		_ = h.tags.SetImageTags(img.ID, tagIDs)
	}

	baseURL := service.ExtractBaseURL(c)
	links := service.BuildLinks(img.ID, img.Ext, baseURL)
	return c.Status(fiber.StatusCreated).JSON(buildImageResponse(img, tagNames, links))
}

// List GET /api/v1/images
func (h *ImageHandler) List(c *fiber.Ctx) error {
	filter := store.ListFilter{
		Tag:   c.Query("tag"),
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 40),
		Sort:  c.Query("sort", "newest"),
	}

	images, total, err := h.images.List(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	baseURL := service.ExtractBaseURL(c)

	items := make([]fiber.Map, 0, len(images))
	for _, img := range images {
		img := img
		tags, _ := h.tags.GetImageTags(img.ID)
		tagNames := make([]string, len(tags))
		for i, t := range tags {
			tagNames[i] = t.Name
		}
		links := service.BuildLinks(img.ID, img.Ext, baseURL)
		items = append(items, buildImageResponse(&img, tagNames, links))
	}

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

// GetOne GET /api/v1/images/:id
func (h *ImageHandler) GetOne(c *fiber.Ctx) error {
	img, err := h.images.GetByID(c.Params("id"))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tags, _ := h.tags.GetImageTags(img.ID)
	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	baseURL := service.ExtractBaseURL(c)
	links := service.BuildLinks(img.ID, img.Ext, baseURL)
	return c.JSON(buildImageResponse(img, tagNames, links))
}

// Delete DELETE /api/v1/images/:id
func (h *ImageHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	img, err := h.images.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.images.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 异步删除磁盘文件
	go func() {
		origPath := filepath.Join(h.dataDir, "images", id+img.Ext)
		optPath := filepath.Join(h.dataDir, "images", id+"_opt.jpg")
		thumbPath := filepath.Join(h.dataDir, "thumbs", id+"_thumb.jpg")
		for _, p := range []string{origPath, optPath, thumbPath} {
			_ = removeFile(p)
		}
		_ = h.tags.DeleteUnused()
	}()

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateTags PUT /api/v1/images/:id/tags
func (h *ImageHandler) UpdateTags(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := h.images.GetByID(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	var body struct {
		Tags []string `json:"tags"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}

	savedTags, err := h.tags.UpsertTags(body.Tags)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tagIDs := make([]int, len(savedTags))
	for i, t := range savedTags {
		tagIDs[i] = t.ID
	}

	if err := h.tags.SetImageTags(id, tagIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tagNames := make([]string, len(savedTags))
	for i, t := range savedTags {
		tagNames[i] = t.Name
	}
	return c.JSON(fiber.Map{"tags": tagNames})
}

func buildImageResponse(img *store.Image, tags []string, links service.Links) fiber.Map {
	return fiber.Map{
		"id":         img.ID,
		"filename":   img.Filename,
		"ext":        img.Ext,
		"size":       img.Size,
		"width":      img.Width,
		"height":     img.Height,
		"mime_type":  img.MimeType,
		"created_at": img.CreatedAt,
		"tags":       tags,
		"url":        links.URL,
		"webp_url":   links.WebPURL,
		"thumb_url":  links.ThumbURL,
		"markdown":   links.Markdown,
		"html":       links.HTML,
		"bbcode":     links.BBCode,
	}
}

func parseTags(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func removeFile(path string) error {
	return os.Remove(path)
}
