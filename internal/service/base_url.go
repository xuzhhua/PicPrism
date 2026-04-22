package service

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ExtractBaseURL 从请求头推断 base URL
// 优先读取 X-Forwarded-Proto + X-Forwarded-Host（Caddy 反代场景）
// 回退到请求本身的 scheme + host
func ExtractBaseURL(c *fiber.Ctx) string {
	proto := c.Get("X-Forwarded-Proto")
	host := c.Get("X-Forwarded-Host")

	if proto == "" {
		if c.Protocol() == "https" {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	if host == "" {
		host = c.Hostname()
	}

	return proto + "://" + strings.TrimRight(host, "/")
}
