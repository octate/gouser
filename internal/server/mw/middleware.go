// Package mw is user Middleware package
package mw

import (
	"gouser/er"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerX(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := c.Errors.Last()
			if err == nil {
				// no errors, abort with success
				return
			}

			log.Error(err.Err.Error())

			e := er.From(err.Err)
			httpStatus := http.StatusInternalServerError
			if e.Status > 0 {
				httpStatus = e.Status
			}
			c.JSON(httpStatus, e)
		}()
		c.Next()
	}
}
