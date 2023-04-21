package server

import (
	"gouser/internal/server/mw"

	"github.com/gin-gonic/gin"
)

func rootRoutes(router *gin.RouterGroup, o *Options) {
	r := router.Group("/")

	// middlewares
	r.Use(mw.ErrorHandlerX(o.Log))

}

func v1Routes(router *gin.RouterGroup, o *Options) {
	r := router.Group("/v1/")

	// middlewares
	r.Use(mw.ErrorHandlerX(o.Log))
	//add new routes here
}
