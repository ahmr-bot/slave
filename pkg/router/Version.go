package router

import (
	"github.com/gin-gonic/gin"
	"os"
)

func HandleVersion(c *gin.Context) {
	// 读取本地文件
	data, err := os.ReadFile("./data/version.txt")
	if err != nil {
		c.String(500, "读取版本文件失败：%v", err)
		return
	}

	// 返回文件内容
	c.String(200, string(data))
}
