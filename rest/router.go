package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/service"
)

func newRouter(svcFactory *service.Factory) *gin.Engine {
	if !logrus.IsLevelEnabled(logrus.DebugLevel) {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		router.Use(gin.Logger())
	}

	c := &Controller{axService: svcFactory.Auxiliary}
	router.Group("/").
		GET("status", c.Status).
		GET("metrics", c.Metrics)

	return router
}
