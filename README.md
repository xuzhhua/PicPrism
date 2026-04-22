# PicPrism

个人自托管图床，支持图片上传、标签管理、多格式链接生成。

## 功能

- 上传 JPG / PNG / GIF / WebP / BMP / TIFF / SVG
- 自动生成 JPEG 优化版 + 缩略图（异步）
- SHA-256 去重
- 标签管理与过滤
- 生成 URL / Markdown / HTML / BBCode 链接
- 暗色 / 亮色主题切换
- 单文件二进制部署（前端 embed 进后端）
- Caddy 反代支持（自动读取 `X-Forwarded-Proto` / `X-Forwarded-Host`）

## 快速开始

### Docker（推荐）

```bash
cp .env.example .env
# 编辑 .env，设置 PICPRISM_TOKEN
docker compose up -d
```

访问 http://localhost:8080

### 本地开发

**前置条件**: Go 1.25+, Node.js 20+

```bash
# 后端
go run ./cmd/picprism

# 前端（另一个终端）
cd web && npm install && npm run dev
```

前端开发服务器运行在 http://localhost:5173，API 请求自动代理到 :8080。

### 构建二进制

```bash
cd web && npm run build && cd ..
go build -o picprism ./cmd/picprism
PICPRISM_TOKEN=your-secret-token ./picprism
```

## 配置

| 环境变量 | 默认值 | 说明 |
|---|---|---|
| `PICPRISM_TOKEN` | （必填） | API 鉴权 Bearer Token |
| `PICPRISM_PORT` | `8080` | 监听端口 |
| `PICPRISM_DATA_DIR` | `/data` | 数据目录（SQLite + 图片） |

## Caddy 反代

```
# Caddyfile
pics.example.com {
    reverse_proxy localhost:8080
}
```

Caddy 会自动传递 `X-Forwarded-Proto` 和 `X-Forwarded-Host`，PicPrism 据此生成正确的外链。

## API

所有写操作需要 `Authorization: Bearer <token>` 请求头。

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/v1/images` | 列出图片（支持 `tag`, `page`, `limit`, `sort`） |
| `POST` | `/api/v1/images` | 上传图片（multipart `file` + 可选 `tags`） |
| `GET` | `/api/v1/images/:id` | 获取单张图片 |
| `DELETE` | `/api/v1/images/:id` | 删除图片 |
| `PUT` | `/api/v1/images/:id/tags` | 更新标签 |
| `GET` | `/api/v1/tags` | 列出所有标签（含数量） |
| `GET` | `/i/:filename` | 访问图片文件（公开） |
| `GET` | `/health` | 健康检查 |

## 数据目录结构

```
/data/
  picprism.db          # SQLite 数据库
  images/              # 原图 + 优化版 JPEG
  thumbs/              # 缩略图
```
