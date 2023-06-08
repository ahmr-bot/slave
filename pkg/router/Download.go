package router

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
)

func HandleDownload(c *gin.Context) {
	filePath := "data" + c.Param("filepath")

	// 获取文件名
	fileName := path.Base(filePath)

	// 1. 设置响应头
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	// 2. 设置 Content-Disposition 响应头，告诉浏览器下载文件
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	// 3. 打开文件
	f, err := os.Open(filePath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	// 4. 将文件内容复制到响应中
	io.Copy(c.Writer, f)
}
