package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/xxxapi/slave/pkg/config"
	"log"
)

func SetMode() {
	// 初始化配置
	conf := config.GetConfig()

	// 模式切换
	if conf.Debug.Debug == true {
		gin.SetMode(gin.DebugMode)
		log.Printf("当前模式：Debug")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Printf("当前模式：Release")
	}
}
