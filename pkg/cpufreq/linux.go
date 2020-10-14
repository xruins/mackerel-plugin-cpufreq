package cpufreq

import (
	"context"
	"fmt"

	proclinux "github.com/c9s/goprocinfo/linux"
	"golang.org/x/xerrors"
)

const (
	pathToProcCPUInfo = "/proc/cpuinfo"
)

type LinuxCPUFreqGetter struct{}

func (l *LinuxCPUFreqGetter) GetCPUFrequencies(ctx context.Context) (map[string]float64, error) {
	cpuInfo, err := proclinux.ReadCPUInfo(pathToProcCPUInfo)
	if err != nil {
		return nil, xerrors.Errorf("failed to get information of CPUs %v", err)
	}

	result := make(map[string]float64, len(cpuInfo.Processors))
	for _, cpu := range cpuInfo.Processors {
		cpuNum := cpu.PhysicalId
		coreNum := cpu.Id
		name := fmt.Sprintf("cpu%dcore%d", cpuNum, coreNum)
		result[name] = cpu.MHz
	}

	return result, nil
}
