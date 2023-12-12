package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/wanliqun/cgo-game-server/proto"
)

const (
	tplRpcSuccessRateMetricKey = "rpc.rate.%s.success"
	tplRpcErrorRateMetricKey   = "rpc.rate.%s.error"
)

var (
	overallRpcRateMetricKey        = "rpc.rate.overall"
	overallRpcSuccessRateMetricKey = rpcSuccessRateMetricKey("overall")
	overallRpcErrorRateMetricKey   = rpcErrorRateMetricKey("overall")

	RPC = newRpcMetrics()
)

type rpcMetrics struct {
	mu         sync.Mutex
	rateTimers map[string]metrics.Timer
}

func newRpcMetrics() *rpcMetrics {
	return &rpcMetrics{
		rateTimers: make(map[string]metrics.Timer),
	}
}

func (m *rpcMetrics) GetOrRegisterTimer(rk string) metrics.Timer {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.rateTimers[rk]
	if !ok {
		t = metrics.GetOrRegisterTimer(rk, nil)
		m.rateTimers[rk] = t
	}

	return t
}

func (m *rpcMetrics) OverallRpcRateTimer() metrics.Timer {
	return m.GetOrRegisterTimer(overallRpcRateMetricKey)
}

func (m *rpcMetrics) IterateRateTimers(cb func(key string, t metrics.Timer)) {
	rateKeys := m.allRateKeys()
	for i := range rateKeys {
		timer := m.GetOrRegisterTimer(rateKeys[i])
		cb(rateKeys[i], timer.Snapshot())
	}
}

func (m *rpcMetrics) allRateKeys() (res []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k := range m.rateTimers {
		res = append(res, k)
	}

	return res
}

func (m *rpcMetrics) Rate(msgType proto.MessageType, err error, start time.Time) {
	metricKeys := []string{overallRpcRateMetricKey}

	if err != nil {
		metricKeys = append(metricKeys, overallRpcErrorRateMetricKey)
		metricKeys = append(metricKeys, rpcErrorRateMetricKey(msgType.String()))
	} else {
		metricKeys = append(metricKeys, overallRpcSuccessRateMetricKey)
		metricKeys = append(metricKeys, rpcSuccessRateMetricKey(msgType.String()))
	}

	for _, mk := range metricKeys {
		m.GetOrRegisterTimer(mk).UpdateSince(start)
	}
}

func rpcSuccessRateMetricKey(index string) string {
	return fmt.Sprintf(tplRpcSuccessRateMetricKey, index)
}

func rpcErrorRateMetricKey(index string) string {
	return fmt.Sprintf(tplRpcSuccessRateMetricKey, index)
}
