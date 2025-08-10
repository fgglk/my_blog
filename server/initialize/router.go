package initialize

import (
	"server/global"
	"server/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)

	return routers.SetupRouter()
}
