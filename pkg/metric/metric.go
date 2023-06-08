package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/net"
	"net/http"
	"sync/atomic"
	"time"
)

type Stats struct {
	IncomingBytes uint64 `json:"incoming_bytes"`
	OutgoingBytes uint64 `json:"outgoing_bytes"`
	RequestCount  uint64 `json:"request_count"`
}

var stats = &Stats{}

func init() {
	go func() {
		for {
			incoming, outgoing := getNetworkTraffic()
			atomic.StoreUint64(&stats.IncomingBytes, incoming)
			atomic.StoreUint64(&stats.OutgoingBytes, outgoing)
			time.Sleep(time.Second)
		}
	}()
}

func getNetworkTraffic() (uint64, uint64) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return 0, 0
	}
	var incoming, outgoing uint64
	for _, c := range counters {
		incoming += c.BytesRecv
		outgoing += c.BytesSent
	}
	return incoming, outgoing
}

func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		atomic.AddUint64(&stats.RequestCount, 1)
		c.Next()
	}
}

func StatsHandler(c *gin.Context) {
	incoming := atomic.LoadUint64(&stats.IncomingBytes)
	outgoing := atomic.LoadUint64(&stats.OutgoingBytes)
	count := atomic.LoadUint64(&stats.RequestCount) / 2

	data := gin.H{
		"incoming_bytes": incoming,
		"outgoing_bytes": outgoing,
		"request_count":  count,
	}

	c.JSON(http.StatusOK, data)
}
