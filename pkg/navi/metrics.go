package navi

import (
	"bufio"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type RoadMetrics struct {
	Traffic      map[string]int64 `json:"traffic"`
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
	if viper.GetBool("performance.low_memory") {
		return
	}

	v.TotalTraffic++
	trace.Timestamp = time.Now()
	if _, ok := v.Traffic[trace.Region]; !ok {
		v.Traffic[trace.Region] = 0
	} else {
		v.Traffic[trace.Region]++
	}

	raw, _ := jsoniter.Marshal(trace)
	accessLogger.Println(string(raw))
}

func (v *RoadMetrics) ReadTrace() []RoadTrace {
	fp := viper.GetString("logging.access")
	file, err := os.Open(fp)
	if err != nil {
		return nil
	}
	defer file.Close()

	var out []RoadTrace
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var entry RoadTrace
		if err := jsoniter.Unmarshal([]byte(line), &entry); err == nil {
			out = append(out, entry)
		}
	}

	return out
}
