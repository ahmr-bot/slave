package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xxxapi/slave/middleware"
	"github.com/xxxapi/slave/pkg/metric"
)

func LoadRoutes(r *gin.Engine) {
	// 加载中间件
	r.Use(metric.StatsMiddleware())
	r.Use(middleware.CORSMiddleware())
	// 加载路由
	r.Use(metric.StatsMiddleware())
	DownloadRouter := r.Group("/download")
	{
		DownloadRouter.GET("/*filepath", HandleDownload)
	}
	APIRouter := r.Group("/api")
	{
		APIRouter.GET("/v1/version", HandleVersion)
		APIRouter.GET("/v1/network", metric.StatsHandler)
	}
	r.GET("/", func(c *gin.Context) {
		c.String(200, "API IS OK!")
	})
}
