package backend

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/duxlabs-dev/dux-go/ipc"
	"github.com/duxlabs-dev/dux-go/log"
)

const sampleInterval = 2 * time.Second

type systemInfo struct {
	Load              float64 `json:"load"`
	MemoryUsedPercent int     `json:"memoryUsedPercent"`
}

func readSystemInfo() (systemInfo, error) {
	loadavg, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return systemInfo{}, err
	}
	meminfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return systemInfo{}, err
	}
	return parseSystemInfo(string(loadavg), string(meminfo))
}

func parseSystemInfo(loadavg, meminfo string) (systemInfo, error) {
	fields := strings.Fields(loadavg)
	if len(fields) == 0 {
		return systemInfo{}, errors.New("loadavg is empty")
	}
	load, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return systemInfo{}, fmt.Errorf("loadavg: %w", err)
	}
	total := meminfoKB(meminfo, "MemTotal")
	if total == 0 {
		return systemInfo{}, errors.New("meminfo has no MemTotal")
	}
	used := total - min(total, meminfoKB(meminfo, "MemAvailable"))
	return systemInfo{
		Load:              load,
		MemoryUsedPercent: int(math.Round(float64(used) * 100 / float64(total))),
	}, nil
}

func meminfoKB(meminfo, key string) uint64 {
	for line := range strings.Lines(meminfo) {
		name, value, found := strings.Cut(line, ":")
		if !found || name != key {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			return 0
		}
		kb, _ := strconv.ParseUint(fields[0], 10, 64)
		return kb
	}
	return 0
}

func (b *Backend) sampleSystemInfo() {
	info, err := readSystemInfo()
	if err != nil {
		log.Warnf("system info: %v", err)
		return
	}
	b.bus.Publish("sysinfo", info)
}

func (b *Backend) publishSystemInfo(ctx context.Context) {
	ticker := time.NewTicker(sampleInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		if !b.bus.HasSubscriber("sysinfo") {
			continue
		}
		b.sampleSystemInfo()
	}
}

func (b *Backend) handleSystemInfoGet(_ context.Context, w *ipc.ConnWriter, req ipc.Request, _ *ipc.Subscriber) {
	info, err := readSystemInfo()
	if err != nil {
		w.RespondError(req.ID, err.Error())
		return
	}
	w.Respond(req.ID, info)
}
