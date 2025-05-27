package router

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xuewenG/cdn-proxy/pkg/config"
	"github.com/xuewenG/cdn-proxy/pkg/handler"
)

func Bind(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(config.Config.CorsOrigin, ","),
		AllowMethods: []string{"GET", "OPTIONS"},
	}))

	r.GET("/:cdnName/*resourcePath", handler.CdnProxy)
}
