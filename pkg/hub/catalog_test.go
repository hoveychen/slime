package hub

import (
	"testing"

	"github.com/hoveychen/slime/pkg/hwinfo"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCatalog_GetHardwareInfo(t *testing.T) {
	mc := NewMemoryCatalog()

	// Test case 1: agentID not found
	hwInfo := mc.GetHardwareInfo(123)
	assert.Nil(t, hwInfo)

	// Test case 2: agentID found
	expectedHWInfo := &hwinfo.HWInfo{CPUCores: 1}
	mc.SetHardwareInfo(123, expectedHWInfo)
	hwInfo = mc.GetHardwareInfo(123)
	assert.Equal(t, expectedHWInfo, hwInfo)
}

func TestMemoryCatalog_SetHardwareInfo(t *testing.T) {
	mc := NewMemoryCatalog()

	// Test case 1: set nil hwInfo
	mc.SetHardwareInfo(123, nil)
	hwInfo := mc.GetHardwareInfo(123)
	assert.Nil(t, hwInfo)

	// Test case 2: set non-nil hwInfo
	expectedHWInfo := &hwinfo.HWInfo{CPUCores: 1}
	mc.SetHardwareInfo(123, expectedHWInfo)
	hwInfo = mc.GetHardwareInfo(123)
	assert.Equal(t, expectedHWInfo, hwInfo)
}

func TestMemoryCatalog_GetPaused(t *testing.T) {
	mc := NewMemoryCatalog()

	// Test case 1: agentID not found
	paused := mc.GetPaused(123)
	assert.False(t, paused)

	// Test case 2: agentID found
	mc.SetPaused(123, true)
	paused = mc.GetPaused(123)
	assert.True(t, paused)
}

func TestMemoryCatalog_SetPaused(t *testing.T) {
	mc := NewMemoryCatalog()

	// Test case 1: set paused to false
	mc.SetPaused(123, false)
	paused := mc.GetPaused(123)
	assert.False(t, paused)

	// Test case 2: set paused to true
	mc.SetPaused(123, true)
	paused = mc.GetPaused(123)
	assert.True(t, paused)
}
