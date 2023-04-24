package server

import (
	"gouser/internal/server/mw"

	"github.com/gin-gonic/gin"
)

func v1Routes(router *gin.RouterGroup, o *Options) {
	r := router.Group("/v1/")

	// middlewares
	r.Use(mw.ErrorHandlerX(o.Log))

	//add new routes here
	r.POST("/users", o.UserHandler.CreateUser)
	r.GET("/users/:user_id", o.UserHandler.FetchUserByID)
	r.GET("/users", o.UserHandler.FetchAllUsers)
	r.PUT("users/:user_id", o.UserHandler.UpdateUser)
}
