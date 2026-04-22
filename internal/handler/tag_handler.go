package handler

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/picprism/picprism/internal/store"
)

type TagHandler struct {
	tags *store.TagStore
}

func NewTagHandler(tags *store.TagStore) *TagHandler {
	return &TagHandler{tags: tags}
}

// ListTags GET /api/v1/tags
func (h *TagHandler) ListTags(c *fiber.Ctx) error {
	tags, err := h.tags.ListAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if tags == nil {
		tags = []store.TagWithCount{}
	}
	return c.JSON(fiber.Map{"tags": tags})
}

// deleteFile removes a file, ignoring not-found errors
func deleteFile(path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		// log silently
		_ = err
	}
}
