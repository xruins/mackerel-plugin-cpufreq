package cpufreq

import (
	"context"
)

type CPUFreqGetter interface {
	GetCPUFrequencies(context.Context) (map[string]float64, error)
}
