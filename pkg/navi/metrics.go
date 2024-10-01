package navi

import (
	"github.com/spf13/viper"
	"time"
)

type RoadMetrics struct {
	Traces []RoadTrace `json:"-"`

	Traffic      map[string]int64 `json:"traffic"`
	TrafficFrom  map[string]int64 `json:"traffic_from"`
	TotalTraffic int64            `json:"total_traffic"`
	StartupAt    time.Time        `json:"startup_at"`
}

type RoadTrace struct {
	Timestamp   time.Time      `json:"timestamp"`
	Region      string         `json:"region"`
	Location    string         `json:"location"`
	Destination string         `json:"destination"`
	Uri         string         `json:"uri"`
	IpAddress   string         `json:"ip_address"`
	UserAgent   string         `json:"user_agent"`
	Error       RoadTraceError `json:"error"`
}

type RoadTraceError struct {
	IsNull  bool   `json:"is_null"`
	Message string `json:"message"`
}

func (v *RoadMetrics) AddTrace(trace RoadTrace) {
	v.TotalTraffic++
	trace.Timestamp = time.Now()
	if _, ok := v.Traffic[trace.Region]; !ok {
		v.Traffic[trace.Region] = 0
	} else {
		v.Traffic[trace.Region]++
	}
	if _, ok := v.TrafficFrom[trace.IpAddress]; !ok {
		v.TrafficFrom[trace.IpAddress] = 0
	} else {
		v.TrafficFrom[trace.IpAddress]++
	}

	v.Traces = append(v.Traces, trace)

	// Garbage recycle
	if len(v.Traffic) > viper.GetInt("performance.traces_limit") {
		v.Traffic = make(map[string]int64)
	}
	if len(v.TrafficFrom) > viper.GetInt("performance.traces_limit") {
		v.TrafficFrom = make(map[string]int64)
	}
	if len(v.Traces) > viper.GetInt("performance.traces_limit") {
		v.Traces = v.Traces[1:]
	}
}
