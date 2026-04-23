package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/picprism/picprism/internal/config"
	"github.com/picprism/picprism/internal/handler"
	"github.com/picprism/picprism/internal/middleware"
	"github.com/picprism/picprism/internal/store"
	webui "github.com/picprism/picprism/web"
)

func main() {
	cfg := config.Load()

	// 确保数据目录存在
	dirs := []string{
		cfg.DataDir,
		filepath.Join(cfg.DataDir, "images"),
		filepath.Join(cfg.DataDir, "thumbs"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("failed to create dir %s: %v", d, err)
		}
	}

	// 初始化数据库
	dbPath := filepath.Join(cfg.DataDir, "picprism.db")
	db, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	imageStore := store.NewImageStore(db)
	tagStore := store.NewTagStore(db)

	imageHandler := handler.NewImageHandler(imageStore, tagStore, cfg.DataDir)
	tagHandler := handler.NewTagHandler(tagStore)
	fileHandler := handler.NewFileHandler(cfg.DataDir)

	app := fiber.New(fiber.Config{
		BodyLimit:      50 * 1024 * 1024, // 50MB
		ReadBufferSize: 16 * 1024,        // 16KB，避免反代追加头后触发 431
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// 图片直链（公开）
	app.Get("/i/:filename", fileHandler.Serve)

	// API 路由
	api := app.Group("/api/v1")
	api.Get("/images", imageHandler.List)
	api.Post("/images", middleware.Auth(cfg.Token), imageHandler.Upload)
	api.Get("/images/:id", imageHandler.GetOne)
	api.Delete("/images/:id", middleware.Auth(cfg.Token), imageHandler.Delete)
	api.Put("/images/:id/tags", middleware.Auth(cfg.Token), imageHandler.UpdateTags)
	api.Get("/tags", tagHandler.ListTags)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// 前端 SPA（嵌入 web/dist）
	distFS, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		log.Fatalf("embed sub: %v", err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("PicPrism starting on %s, data dir: %s", addr, cfg.DataDir)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
