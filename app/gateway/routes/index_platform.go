package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzug/goS/app/gateway/http"
)

func IndexPlatformRegisterHandlers(rg *gin.RouterGroup) {
	indexPlatformGroup := rg.Group("/index_platform")
	{
		indexPlatformGroup.POST("/build_index", http.BuildIndexByFiles)
	}
}
