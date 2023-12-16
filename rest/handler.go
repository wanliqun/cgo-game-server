package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wanliqun/cgo-game-server/service"
)

type Controller struct {
	axService *service.AuxiliaryService
}

type ServerStatus struct {
	*service.ServerStatus
	Uptime string
}

func (c *Controller) Status(ctx *gin.Context) {
	srvStat := c.axService.CollectServerStatus()
	ctx.JSON(http.StatusOK, &ServerStatus{
		ServerStatus: srvStat,
		Uptime:       srvStat.Uptime.String(),
	})
}

func (c *Controller) Metrics(ctx *gin.Context) {
	metrics := c.axService.GatherAllRPCRateMetrics()
	ctx.JSON(http.StatusOK, metrics)
}
