package main

import (
	"context"
	"fmt"
	"runtime"

	"flag"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/xruins/mackerel-plugin-cpufreq/pkg/cpufreq"
	"golang.org/x/xerrors"
)

type CPUFreqPlugin struct {
	prefix string
}

func (htp *CPUFreqPlugin) FetchMetrics() (map[string]float64, error) {
	ctx := context.Background()
	var getter cpufreq.CPUFreqGetter

	switch runtime.GOOS {
	case "linux":
		getter = &cpufreq.LinuxCPUFreqGetter{}
	default:
		return nil, xerrors.New("unsupported platform")
	}

	result, err := getter.GetCPUFrequencies(ctx)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]float64, len(result))
	for k, v := range result {
		key := fmt.Sprintf("%s.frequency", k)
		metrics[key] = v
	}

	return metrics, nil
}
func (htp *CPUFreqPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func (htp *CPUFreqPlugin) MetricKeyPrefix() string {
	if htp.prefix == "" {
		return "cpufreq"
	}
	return htp.prefix
}

var graphdef = map[string]mp.Graphs{
	"#": {
		Label: "CPU Frequency",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "frequency", Label: "Frequency", Diff: false},
		},
	},
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	cpufreq := CPUFreqPlugin{
		prefix: *optPrefix,
	}

	plugin := mp.NewMackerelPlugin(&cpufreq)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
