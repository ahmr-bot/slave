package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 为Gin绑定CORS跨域访问处理中间件
func CORSMiddleware() gin.HandlerFunc {
	// 初始化配置
	return cors.Default()
}
