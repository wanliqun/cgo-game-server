package metrics

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/wanliqun/cgo-game-server/proto"
)

const (
	onlinePlayersMetricKey = "concurrent_players"
	udpConnsMetricKey      = "udp_connections"
	tcpConnsMetricKey      = "tcp_connections"

	tplRpcSuccessRateMetricKey = "rpc.rate.%s.success"
	tplRpcErrorRateMetricKey   = "rpc.rate.%s.error"
)

var (
	overallRpcSuccessRateMetricKey = rpcSuccessRateMetricKey("overall")
	overallRpcErrorRateMetricKey   = rpcErrorRateMetricKey("overall")

	Server = &serverMetrics{}
	RPC    = &rpcMetrics{}
)

type serverMetrics struct{}

func (m *serverMetrics) OnlinePlayers() metrics.Gauge {
	return metrics.GetOrRegisterGauge(onlinePlayersMetricKey, nil)
}

func (m *serverMetrics) UDPConnections() metrics.Gauge {
	return metrics.GetOrRegisterGauge(udpConnsMetricKey, nil)
}

func (m *serverMetrics) TCPConnections() metrics.Gauge {
	return metrics.GetOrRegisterGauge(tcpConnsMetricKey, nil)
}

type rpcMetrics struct{}

func (m *rpcMetrics) Rate(msgType proto.MessageType, err error, start time.Time) {
	if err != nil {
		metrics.GetOrRegisterTimer(overallRpcErrorRateMetricKey, nil).UpdateSince(start)
		t := metrics.GetOrRegisterTimer(rpcErrorRateMetricKey(msgType.String()), nil)
		t.UpdateSince(start)
		return
	}

	metrics.GetOrRegisterTimer(overallRpcSuccessRateMetricKey, nil).UpdateSince(start)
	t := metrics.GetOrRegisterTimer(rpcSuccessRateMetricKey(msgType.String()), nil)
	t.UpdateSince(start)
}

func rpcSuccessRateMetricKey(index string) string {
	return fmt.Sprintf(tplRpcSuccessRateMetricKey, index)
}

func rpcErrorRateMetricKey(index string) string {
	return fmt.Sprintf(tplRpcSuccessRateMetricKey, index)
}
