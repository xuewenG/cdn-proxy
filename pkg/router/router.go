package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xuewenG/cdn-proxy/pkg/handler"
)

func Bind(r *gin.Engine) {
	r.GET("/:cdnName/*resourcePath", handler.CdnProxy)
}
