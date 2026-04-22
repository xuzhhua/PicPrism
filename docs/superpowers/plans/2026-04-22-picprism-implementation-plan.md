# PicPrism — 实现计划

**关联规格：** `docs/superpowers/specs/2026-04-22-picprism-design.md`  
**日期：** 2026-04-22

---

## 阶段概览

| 阶段 | 内容 | 产出 |
|------|------|------|
| 1 | 项目骨架 & 配置 | Go module、目录结构、环境变量 |
| 2 | 数据库层 | SQLite schema、store 接口 |
| 3 | 图片处理服务 | 上传、WebP 转换、缩略图 |
| 4 | REST API | 全部端点、Token 鉴权中间件 |
| 5 | 前端基础 | Vue 3 项目、路由、Pinia、主题 |
| 6 | 前端界面 | 三栏布局、图库、详情面板、上传 |
| 7 | 嵌入 & 构建 | embed.FS、Dockerfile、docker-compose |
| 8 | 集成收尾 | 端到端验证、Caddy 配置文档 |

---

## 阶段 1：项目骨架 & 配置

### 任务
1. 初始化 Go module：`go mod init github.com/picprism/picprism`
2. 创建目录结构（`cmd/picprism`、`internal/{config,handler,service,store,middleware}`、`web/`）
3. 实现 `internal/config/config.go`：读取 `PICPRISM_TOKEN`、`PICPRISM_PORT`（默认 8080）、`PICPRISM_DATA_DIR`（默认 `/data`）
4. 创建 `cmd/picprism/main.go`：加载 config，创建数据目录，启动 Fiber 服务器（占位路由）
5. 创建 `.env.example`

### 验收
- `go build ./cmd/picprism` 成功
- 服务器监听指定端口，`/health` 返回 200

---

## 阶段 2：数据库层

### 任务
1. 添加依赖：`modernc.org/sqlite`、`github.com/jmoiron/sqlx`
2. 实现 `internal/store/db.go`：初始化 SQLite，执行建表 SQL（images、tags、image_tags）
3. 实现 `internal/store/image_store.go`：
   - `Create(image Image) error`
   - `GetByID(id string) (*Image, error)`
   - `GetByHash(hash string) (*Image, error)`
   - `List(filter ListFilter) ([]Image, int, error)`
   - `Delete(id string) error`
4. 实现 `internal/store/tag_store.go`：
   - `UpsertTags(names []string) ([]Tag, error)`
   - `SetImageTags(imageID string, tagIDs []int) error`
   - `GetImageTags(imageID string) ([]Tag, error)`
   - `ListAll() ([]TagWithCount, error)`

### 验收
- 单元测试：`store_test.go` 覆盖 Create / GetByID / List / Delete / Tag 操作
- 测试使用内存 SQLite（`:memory:`）

---

## 阶段 3：图片处理服务

### 任务
1. 添加依赖：`github.com/jaevor/go-nanoid`、`golang.org/x/image`（WebP 编码）
2. 实现 `internal/service/image_service.go`：
   - `ProcessUpload(file io.Reader, originalName string) (*ProcessResult, error)`
     - 计算 SHA-256
     - 检测格式、解码、读取宽高
     - 保存原图至 `<dataDir>/images/<id>.<ext>`
     - 启动 goroutine：生成 WebP → `<dataDir>/images/<id>.webp`
     - 启动 goroutine：生成缩略图（800px 限宽）→ `<dataDir>/thumbs/<id>_thumb.webp`
   - `BuildLinks(id, ext, baseURL string) Links`
     - 返回 url、webp_url、thumb_url、markdown、html、bbcode
3. 实现 `internal/service/base_url.go`：
   - `ExtractBaseURL(r *fiber.Ctx) string`：读取 `X-Forwarded-Proto` + `X-Forwarded-Host`，fallback 到请求本身的 host

### 支持格式处理矩阵

| 格式 | 解码 | WebP 转换 | 缩略图 |
|------|------|-----------|--------|
| JPEG | `image/jpeg` | ✓ | ✓ |
| PNG | `image/png` | ✓ | ✓ |
| GIF | `image/gif`（静态帧）| ✓ | ✓ |
| WebP | `golang.org/x/image/webp` | 直接复用 | ✓ |
| BMP | `golang.org/x/image/bmp` | ✓ | ✓ |
| TIFF | `golang.org/x/image/tiff` | ✓ | ✓ |
| AVIF | `github.com/gen2brain/avif` | ✓ | ✓ |
| SVG | 不解码，直存 | 不转换，webp_url = 原图 | 不生成 |

### 验收
- 单元测试：上传各格式测试文件，验证磁盘文件生成、宽高读取正确
- WebP 缩略图宽度 ≤ 800px

---

## 阶段 4：REST API & 鉴权

### 任务
1. 实现 `internal/middleware/auth.go`：
   - 检查 `Authorization: Bearer <token>`，不匹配返回 401
2. 实现各 handler（每个文件一个资源）：
   - `handler/image_handler.go`：`Upload`、`List`、`GetOne`、`Delete`
   - `handler/tag_handler.go`：`ListTags`、`UpdateImageTags`
   - `handler/file_handler.go`：`ServeFile`（原图 + WebP + 缩略图直链）
3. 在 `main.go` 注册所有路由，写操作路由组套用 auth 中间件
4. 统一错误响应格式：`{"error": "message"}`

### 路由注册示意

```go
api := app.Group("/api/v1")
api.Get("/images", h.List)
api.Post("/images", authMiddleware, h.Upload)
api.Get("/images/:id", h.GetOne)
api.Delete("/images/:id", authMiddleware, h.Delete)
api.Put("/images/:id/tags", authMiddleware, h.UpdateTags)
api.Get("/tags", h.ListTags)

app.Get("/i/:filename", fileHandler.Serve)
```

### 验收
- 使用 `httptest` 测试：未授权 POST 返回 401，合法 token 返回 201
- 上传后能通过 GET `/i/:id.webp` 取回图片

---

## 阶段 5：前端基础

### 任务
1. 在 `web/` 创建 Vue 3 + Vite + TypeScript 项目（`npm create vite@latest`）
2. 安装依赖：`tailwindcss`、`pinia`、`vue-router`、`@vueuse/core`
3. 配置 Tailwind（自定义 CSS 变量驱动的颜色 token）
4. 实现 `stores/theme.ts`：亮/暗模式切换，持久化到 `localStorage`，默认跟随系统
5. 实现 `stores/auth.ts`：存取 API token（`localStorage`），提供登录/登出
6. 实现 `api/client.ts`：封装 fetch，自动注入 token header，处理 401
7. 配置 Vite proxy：开发模式下 `/api` 和 `/i` 代理到 `localhost:8080`

### 验收
- `npm run dev` 启动，主题切换正常，token 输入框显示/隐藏逻辑正确

---

## 阶段 6：前端界面

### 任务

#### 6.1 布局骨架
- `components/AppNav.vue`：顶部导航（Logo + 上传按钮 + 主题切换）
- `components/AppSidebar.vue`：标签导航列表
- `views/GalleryView.vue`：三栏主视图容器

#### 6.2 图库
- `components/ImageGrid.vue`：网格视图，懒加载缩略图
- `components/ImageList.vue`：列表视图（文件名、尺寸、标签、时间）
- `components/ImageCard.vue`：卡片，悬浮显示操作按钮（查看详情 / 复制 WebP 链接 / 删除）
- `stores/gallery.ts`：分页拉取、标签过滤、搜索、排序状态

#### 6.3 详情面板
- `components/DetailPanel.vue`：
  - 元数据展示
  - 标签 inline 编辑（添加/删除）
  - 链接区（直链 / Markdown / HTML / BBCode），点击复制，复制后提示

#### 6.4 上传
- `components/UploadZone.vue`：拖拽 + 点击选文件，支持多文件，上传进度条
- 上传成功后自动刷新图库并选中新图片

#### 6.5 搜索 & 过滤
- 工具栏搜索框：实时 debounce 过滤（300ms）
- 侧边栏标签点击：单选过滤

### 验收
- 网格/列表视图切换正常
- 上传 → 显示 → 复制链接 → 删除 完整流程可走通
- 主题切换后界面无闪烁

---

## 阶段 7：嵌入 & 构建

### 任务
1. 在 `cmd/picprism/main.go` 添加 `//go:embed web/dist` 指令
2. 配置 Fiber 静态文件服务从 `embed.FS` 读取，SPA fallback 到 `index.html`
3. 编写 `Dockerfile`：
   - Stage 1（node:20-alpine）：构建前端 `npm run build`
   - Stage 2（golang:1.23-alpine）：构建 Go 二进制
   - Stage 3（alpine:3.20）：最终镜像，只含二进制
4. 编写 `docker-compose.yml`
5. 创建 `.env.example`
6. 在 `Makefile` 添加：`build`、`docker-build`、`dev`（同时启动后端 + 前端 dev server）目标

### 验收
- `docker build` 成功，镜像大小 < 30MB
- 容器启动后访问 `:8080` 显示前端界面

---

## 阶段 8：集成收尾

### 任务
1. 端到端手动测试：
   - 上传 JPEG / PNG / GIF / WebP / SVG，验证链接正确、WebP 可访问
   - Caddy 反代后 `X-Forwarded-Host` 推断 Base URL 正确
   - Token 错误返回 401
   - SHA-256 去重：同一文件上传两次返回已有记录
2. 编写 `README.md`：快速启动、Caddy 配置示例、环境变量说明
3. 在项目根创建 `Caddyfile.example`

### 验收
- 全部阶段验收项通过
- README 可独立指导完成部署

---

## 依赖清单（Go）

```
github.com/gofiber/fiber/v2
modernc.org/sqlite
github.com/jmoiron/sqlx
github.com/jaevor/go-nanoid
golang.org/x/image
github.com/gen2brain/avif
```

## 依赖清单（前端）

```
vue@3, vite, typescript
tailwindcss, @tailwindcss/vite
pinia
vue-router
@vueuse/core
```
