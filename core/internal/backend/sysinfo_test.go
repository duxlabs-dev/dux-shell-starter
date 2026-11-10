package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const meminfoSample = `MemTotal:       32000000 kB
MemFree:         2000000 kB
MemAvailable:   24000000 kB
Buffers:          500000 kB
`

func TestParseSystemInfo(t *testing.T) {
	info, err := parseSystemInfo("0.52 0.41 0.30 1/1234 56789\n", meminfoSample)
	require.NoError(t, err)
	assert.Equal(t, systemInfo{Load: 0.52, MemoryUsedPercent: 25}, info)
}

func TestParseSystemInfoRejectsMissingFields(t *testing.T) {
	tests := []struct {
		name    string
		loadavg string
		meminfo string
	}{
		{"empty loadavg", "", meminfoSample},
		{"loadavg not a number", "busy 0.41 0.30", meminfoSample},
		{"no MemTotal", "0.52 0.41 0.30", "MemAvailable: 24000000 kB\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseSystemInfo(tc.loadavg, tc.meminfo)
			assert.Error(t, err)
		})
	}
}

func TestReadSystemInfoOnThisHost(t *testing.T) {
	info, err := readSystemInfo()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, info.MemoryUsedPercent, 0)
	assert.LessOrEqual(t, info.MemoryUsedPercent, 100)
}
