# PicPrism — 设计规格文档

**日期：** 2026-04-22  
**状态：** 已批准

---

## 1. 概述

PicPrism 是一个个人自托管图床应用，提供图片上传、存储、检索、浏览、管理及外链生成功能。通过 Caddy 反代对外提供服务，使用 API Token 控制写权限，图片浏览无需鉴权。

---

## 2. 技术栈

| 层次 | 技术 |
|------|------|
| 后端 | Go + Fiber |
| 前端 | Vue 3 + Vite + Tailwind CSS + Pinia |
| 数据库 | SQLite（`modernc.org/sqlite`，纯 Go，无 CGO）|
| 图片处理 | Go 标准库 `image/*` + `golang.org/x/image`（WebP 编码）|
| 部署 | Docker + Docker Compose |
| 反代 | Caddy（自动 TLS，注入转发头）|

---

## 3. 架构

```
Caddy (TLS 终止 + 反代)
    │  X-Forwarded-Host / X-Forwarded-Proto
    ▼
PicPrism 后端 (Go + Fiber, :8080)
    ├── REST API  /api/v1/*   (Token 鉴权)
    ├── 图片直链  /i/:id.*    (公开)
    └── 前端 SPA  /*          (embed.FS 嵌入)
         ├── SQLite  /data/picprism.db
         ├── 原图    /data/images/
         └── 缩略图  /data/thumbs/
```

**关键决定：**
- 前端构建产物通过 `embed.FS` 嵌入 Go 二进制，部署只需一个可执行文件 + `/data` 目录。
- Base URL 从 `X-Forwarded-Proto` + `X-Forwarded-Host` 请求头动态推断，无需环境变量配置。
- 图片 ID 使用 nanoid（21 字符），同时作为磁盘文件名。

---

## 4. 数据模型

```sql
CREATE TABLE images (
    id         TEXT PRIMARY KEY,   -- nanoid，文件名前缀
    filename   TEXT NOT NULL,      -- 原始文件名
    ext        TEXT NOT NULL,      -- 原始扩展名（.jpg/.png/…）
    size       INTEGER NOT NULL,   -- 字节数
    width      INTEGER NOT NULL,
    height     INTEGER NOT NULL,
    mime_type  TEXT NOT NULL,
    hash       TEXT UNIQUE,        -- SHA-256，防重复上传
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE image_tags (
    image_id TEXT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    tag_id   INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (image_id, tag_id)
);
```

---

## 5. API 端点

所有 `/api/v1/*` 写操作需要 Header：`Authorization: Bearer <token>`

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| `POST` | `/api/v1/images` | 需要 | 上传图片（multipart/form-data，字段 `file`，可选 `tags`）|
| `GET` | `/api/v1/images` | 无 | 列表，查询参数：`tag`、`page`（默认1）、`limit`（默认40）、`sort`（`newest`/`oldest`/`name`/`size`）|
| `GET` | `/api/v1/images/:id` | 无 | 图片元数据 + 链接 |
| `DELETE` | `/api/v1/images/:id` | 需要 | 删除图片及磁盘文件 |
| `PUT` | `/api/v1/images/:id/tags` | 需要 | 全量更新图片标签（`{"tags":["a","b"]}`）|
| `GET` | `/api/v1/tags` | 无 | 所有标签及图片数量 |
| `GET` | `/i/:id.:ext` | 无 | 原图直链（`ext` 匹配原始格式）|
| `GET` | `/i/:id.webp` | 无 | WebP 转换版本 |
| `GET` | `/i/:id_thumb.webp` | 无 | WebP 缩略图（800px 限宽）|

### 上传响应示例

```json
{
  "id": "V1StGXR8_Z5jdHi6B-myT",
  "url":      "https://img.example.com/i/V1StGXR8_Z5jdHi6B-myT.jpg",
  "webp_url": "https://img.example.com/i/V1StGXR8_Z5jdHi6B-myT.webp",
  "thumb_url":"https://img.example.com/i/V1StGXR8_Z5jdHi6B-myT_thumb.webp",
  "markdown": "![alt](https://img.example.com/i/V1StGXR8_Z5jdHi6B-myT.webp)",
  "html":     "<img src=\"https://img.example.com/i/V1StGXR8_Z5jdHi6B-myT.webp\" alt=\"\">",
  "bbcode":   "[img]https://img.example.com/i/V1StGXR8_Z5jdHi6B-myT.webp[/img]"
}
```

---

## 6. 图片处理

**上传流程：**
1. 接收文件，计算 SHA-256，查重（已存在则返回现有记录）
2. 解码图片，读取宽高、MIME 类型
3. 将原图保存至 `/data/images/<id>.<ext>`
4. 异步生成 WebP 版本 → `/data/images/<id>.webp`
5. 异步生成缩略图（800px 限宽，质量 80）→ `/data/thumbs/<id>_thumb.webp`
6. 写入 SQLite 元数据

**支持格式：** JPEG · PNG · GIF（静态帧）· WebP · BMP · TIFF · AVIF · SVG  
**SVG 特殊处理：** 直接存储，不做光栅化转换；WebP 链接回退到 SVG 原文件

---

## 7. 前端界面

**布局：** 三栏式
- 左：标签导航栏（全部图片、最近上传、标签列表）
- 中：工具栏（搜索 + 排序 + 视图切换）+ 网格/列表图库
- 右：详情面板（元数据 + 标签编辑 + 链接复制区）

**主题：** 纯黑白灰配色，无彩色强调色；亮色模式背景 `#fafafa`，暗色 `#0a0a0a`；切换按钮常驻导航栏

**交互细节：**
- 悬浮图片卡片显示操作按钮（查看详情 / 复制 WebP 链接 / 删除）
- 详情面板中各链接格式（直链 / Markdown / HTML / BBCode）点击一键复制
- 支持拖拽上传（导航栏按钮 + 网格区顶部拖拽区域）
- 标签支持多选过滤，可在详情面板内联添加/删除

---

## 8. 鉴权

- 配置项 `PICPRISM_TOKEN`（环境变量）
- 写操作 Header：`Authorization: Bearer <token>`
- 图片直链及列表读取：公开，无需 token
- 前端管理界面：本地 localStorage 存储 token，首次访问若 token 未配置则显示输入框

---

## 9. 部署

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PICPRISM_TOKEN` | （必填）| API 鉴权 token |
| `PICPRISM_PORT` | `8080` | 监听端口 |
| `PICPRISM_DATA_DIR` | `/data` | 数据根目录 |

### docker-compose.yml

```yaml
services:
  picprism:
    image: picprism:latest
    restart: unless-stopped
    volumes:
      - ./data:/data
    environment:
      - PICPRISM_TOKEN=your_secret_token_here
      - PICPRISM_PORT=8080
    ports:
      - "127.0.0.1:8080:8080"
```

### Caddyfile

```caddyfile
img.example.com {
    reverse_proxy localhost:8080
}
```

Caddy 自动注入 `X-Forwarded-Host` / `X-Forwarded-Proto`，PicPrism 读取后动态拼接外链，无需额外配置 base URL。

---

## 10. 项目目录结构

```
PicPrism/
├── cmd/picprism/main.go
├── internal/
│   ├── config/       # 环境变量读取
│   ├── handler/      # HTTP 处理器
│   ├── service/      # 图片处理、链接生成
│   ├── store/        # SQLite 数据访问层
│   └── middleware/   # Token 鉴权
├── web/              # Vue 3 前端
│   └── src/
│       ├── components/
│       ├── views/
│       └── stores/
├── data/             # 运行时（gitignore）
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

---

## 11. 不在范围内（本版本）

- 多用户 / 多 token
- 外部对象存储（S3 / MinIO）
- 图片 EXIF 数据提取
- 图片压缩率对比展示
- 批量操作（批量打标签、批量删除）
