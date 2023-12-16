package service

import (
	"fmt"
	"time"

	gometrics "github.com/rcrowley/go-metrics"
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/metrics"
	"github.com/wanliqun/cgo-game-server/server"
)

type ServerStatus struct {
	ServerName       string
	Uptime           time.Duration
	NumOnlinePlayers int32
	TotalConnections int32
}

type AuxiliaryService struct {
	common.MonickerGenerator
	Config    *config.Config
	playerSvc *PlayerService
	sessMgr   *server.SessionManager
	start     time.Time
}

func NewAuxiliaryService(
	cfg *config.Config, g common.MonickerGenerator,
	svc *PlayerService, mgr *server.SessionManager) *AuxiliaryService {
	return &AuxiliaryService{
		MonickerGenerator: g,
		Config:            cfg,
		playerSvc:         svc,
		sessMgr:           mgr,
		start:             time.Now(),
	}
}

func (s *AuxiliaryService) CollectServerStatus() *ServerStatus {
	return &ServerStatus{
		ServerName:       s.Config.Server.Name,
		Uptime:           time.Since(s.start),
		NumOnlinePlayers: int32(s.playerSvc.Count()),
		TotalConnections: int32(s.sessMgr.Count()),
	}
}

func (s *AuxiliaryService) GatherOverallRPCRateMetrics() map[string]string {
	t := metrics.RPC.OverallRpcRateTimer()
	return map[string]string{
		"requests":       fmt.Sprintf("%d", t.Count()),
		"TPS (m1)":       fmt.Sprintf("%.1f", t.Rate1()),
		"TPS (m5)":       fmt.Sprintf("%.1f", t.Rate5()),
		"TPS (m15)":      fmt.Sprintf("%.1f", t.Rate15()),
		"Latency (mean)": fmt.Sprintf("%.1fms", t.Mean()/1e6),
		"Latency (p75)":  fmt.Sprintf("%.1fms", t.Percentile(0.75)/1e6),
		"Latency (p90)":  fmt.Sprintf("%.1fms", t.Percentile(0.90)/1e6),
		"Latency (p99)":  fmt.Sprintf("%.1fms", t.Percentile(0.99)/1e6),
	}
}

func (s *AuxiliaryService) GatherAllRPCRateMetrics() map[string]string {
	rpcRateMetrics := make(map[string]string)
	metrics.RPC.IterateRateTimers(func(key string, t gometrics.Timer) {
		// Samples count
		sampleCount := fmt.Sprintf("%s Sample Count", key)
		rpcRateMetrics[sampleCount] = fmt.Sprintf("%v", t.Count())

		// TPS
		m1Tps := fmt.Sprintf("%s m1 TPS", key)
		rpcRateMetrics[m1Tps] = fmt.Sprintf("%.1f", t.Rate1())
		m5Tps := fmt.Sprintf("%s m5 TPS", key)
		rpcRateMetrics[m5Tps] = fmt.Sprintf("%.1f", t.Rate5())
		m15Tps := fmt.Sprintf("%s m15 TPS", key)
		rpcRateMetrics[m15Tps] = fmt.Sprintf("%.1f", t.Rate15())

		// Latency
		minLatency := fmt.Sprintf("%s Min Latency", key)
		rpcRateMetrics[minLatency] = fmt.Sprintf("%.1f(ms)", float64(t.Min())/1e6)
		meanLatency := fmt.Sprintf("%s Mean Latency", key)
		rpcRateMetrics[meanLatency] = fmt.Sprintf("%.1f(ms)", t.Mean()/1e6)
		maxLatency := fmt.Sprintf("%s Max Latency", key)
		rpcRateMetrics[maxLatency] = fmt.Sprintf("%1f(ms)", float64(t.Max())/1e6)

		p50Latency := fmt.Sprintf("%s p50 Latency", key)
		rpcRateMetrics[p50Latency] = fmt.Sprintf("%.1f(ms)", t.Percentile(0.5)/1e6)
		p75Latency := fmt.Sprintf("%s p75 Latency", key)
		rpcRateMetrics[p75Latency] = fmt.Sprintf("%.1f(ms)", t.Percentile(0.75)/1e6)
		p90Latency := fmt.Sprintf("%s p90 Latency", key)
		rpcRateMetrics[p90Latency] = fmt.Sprintf("%.1f(ms)", t.Percentile(0.90)/1e6)
		p99Latency := fmt.Sprintf("%s p99 Latency", key)
		rpcRateMetrics[p99Latency] = fmt.Sprintf("%.1f(ms)", t.Percentile(0.99)/1e6)
	})

	return rpcRateMetrics
}
