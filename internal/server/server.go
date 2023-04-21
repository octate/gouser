package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module invokes mainserver
var Module = fx.Options(
	fx.Invoke(
		Run,
	),
)

const (
	addr = "0.0.0.0"
)

// Options is function arguments struct of `Run` function.
type Options struct {
	fx.In

	Config *viper.Viper
	Log    *zap.Logger

	PostgresDB *sql.DB `name:"userDB"`
}

// Run starts the mainserver REST API server
func Run(o Options) {
	router := SetupRouter(&o)
	router.Run(fmt.Sprintf("%s:%s", addr, o.Config.GetString("port")))

	return
}

// SetupRouter creates gin router and registers all user routes to it
func SetupRouter(o *Options) (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.Ginzap(o.Log, time.RFC3339, false))

	// Logs all panic to error log
	// stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(o.Log, true))

	// router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
	// 	SkipPaths: []string{"/_healthz", "/_readyz"},
	// }))

	// Health routes. DO NOT MOVE IT FROM HERE!
	router.GET("/_healthz", HealthHandler(o))
	router.GET("/_readyz", HealthHandler(o))

	rootRouter := router.Group("/")
	rootRoutes(rootRouter, o)

	v1Routes(rootRouter, o)

	return
}

// HealthHandler
func HealthHandler(o *Options) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		err = o.PostgresDB.Ping()
		if err != nil {
			c.AbortWithError(http.StatusFailedDependency, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": "ok"})
	}
}
