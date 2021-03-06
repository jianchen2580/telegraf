package system

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/require"
)

func TestMemStats(t *testing.T) {
	var mps MockPS
	var err error
	defer mps.AssertExpectations(t)
	var acc testutil.Accumulator

	vms := &mem.VirtualMemoryStat{
		Total:     12400,
		Available: 7600,
		Used:      5000,
		Free:      1235,
		Active:    8134,
		Inactive:  1124,
		Slab:      1234,
		Wired:     134,
		// Buffers:     771,
		// Cached:      4312,
		// Shared:      2142,
	}

	mps.On("VMStat").Return(vms, nil)

	err = (&MemStats{&mps}).Gather(&acc)
	require.NoError(t, err)

	memfields := map[string]interface{}{
		"total":             uint64(12400),
		"available":         uint64(7600),
		"used":              uint64(5000),
		"available_percent": float64(7600) / float64(12400) * 100,
		"used_percent":      float64(5000) / float64(12400) * 100,
		"free":              uint64(1235),
		"cached":            uint64(0),
		"buffered":          uint64(0),
		"active":            uint64(8134),
		"inactive":          uint64(1124),
		"wired":             uint64(134),
		"slab":              uint64(1234),
	}
	acc.AssertContainsTaggedFields(t, "mem", memfields, make(map[string]string))

	acc.Metrics = nil
}
