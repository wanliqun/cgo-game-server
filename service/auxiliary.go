package service

import (
	"fmt"

	gometrics "github.com/rcrowley/go-metrics"
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/metrics"
)

type ServerStatus struct {
	NumOnlinePlayers  int64
	NumTCPConnections int64
	NumUDPConnections int64
}

type AuxiliaryService struct {
	common.MonickerGenerator

	Config *config.Config
}

func NewAuxiliaryService(
	cfg *config.Config, generator common.MonickerGenerator) *AuxiliaryService {
	return &AuxiliaryService{Config: cfg, MonickerGenerator: generator}
}

func (s *AuxiliaryService) CollectServerStatus() *ServerStatus {
	return &ServerStatus{
		NumOnlinePlayers:  metrics.Server.OnlinePlayers().Value(),
		NumTCPConnections: metrics.Server.TCPConnections().Value(),
		NumUDPConnections: metrics.Server.UDPConnections().Value(),
	}
}

func (s *AuxiliaryService) GatherRPCRateMetrics() map[string]string {
	rpcRateMetrics := make(map[string]string)
	metrics.RPC.IterateRateTimers(func(key string, t gometrics.Timer) {
		// Samples count
		sampleCount := fmt.Sprintf("%s Sample Count", key)
		rpcRateMetrics[sampleCount] = fmt.Sprintf("%v", t.Count())

		// TPS
		m1Tps := fmt.Sprintf("%s m1 TPS", key)
		rpcRateMetrics[m1Tps] = fmt.Sprintf("%.2f", t.Rate1())
		m5Tps := fmt.Sprintf("%s m5 TPS", key)
		rpcRateMetrics[m5Tps] = fmt.Sprintf("%.2f", t.Rate5())
		m15Tps := fmt.Sprintf("%s m15 TPS", key)
		rpcRateMetrics[m15Tps] = fmt.Sprintf("%.2f", t.Rate15())

		// Latency
		minLatency := fmt.Sprintf("%s Min Latency", key)
		rpcRateMetrics[minLatency] = fmt.Sprintf("%v", t.Min())
		meanLatency := fmt.Sprintf("%s Mean Latency", key)
		rpcRateMetrics[meanLatency] = fmt.Sprintf("%.2f", t.Mean())
		maxLatency := fmt.Sprintf("%s Max Latency", key)
		rpcRateMetrics[maxLatency] = fmt.Sprintf("%v", t.Max())

		p50Latency := fmt.Sprintf("%s p50 Latency", key)
		rpcRateMetrics[p50Latency] = fmt.Sprintf("%.2f", t.Percentile(50))
		p75Latency := fmt.Sprintf("%s p75 Latency", key)
		rpcRateMetrics[p75Latency] = fmt.Sprintf("%.2f", t.Percentile(75))
		p90Latency := fmt.Sprintf("%s p90 Latency", key)
		rpcRateMetrics[p90Latency] = fmt.Sprintf("%.2f", t.Percentile(90))
		p99Latency := fmt.Sprintf("%s p99 Latency", key)
		rpcRateMetrics[p99Latency] = fmt.Sprintf("%.2f", t.Percentile(99))
	})

	return rpcRateMetrics
}
