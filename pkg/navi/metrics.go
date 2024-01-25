package navi

import "github.com/spf13/viper"

type RoadTrace struct {
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

func (v *RoadApp) AddTrace(trace RoadTrace) {
	v.Traces = append(v.Traces, trace)
	if len(v.Traces) > viper.GetInt("performance.traces_limit") {
		v.Traces = v.Traces[1:]
	}
}
