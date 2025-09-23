package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	if cfip := c.GetHeader("CF-Connecting-IP"); cfip != "" {
		return cfip
	}

	if xf := c.GetHeader("X-Forwarded"); xf != "" {
		return strings.TrimSpace(strings.Split(xf, ",")[0])
	}

	return c.ClientIP()
}
