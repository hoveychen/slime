package hub

import (
	"sync"

	"github.com/hoveychen/slime/pkg/hwinfo"
)

type Catalog interface {
	GetHardwareInfo(agentID int) *hwinfo.HWInfo
	SetHardwareInfo(agentID int, hwInfo *hwinfo.HWInfo)
}

type MemoryCatalog struct {
	hwInfoMap map[int]*hwinfo.HWInfo
	mutex     sync.RWMutex
}

func NewMemoryCatalog() *MemoryCatalog {
	return &MemoryCatalog{
		hwInfoMap: make(map[int]*hwinfo.HWInfo),
	}
}

func (mc *MemoryCatalog) GetHardwareInfo(agentID int) *hwinfo.HWInfo {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.hwInfoMap[agentID]
}

func (mc *MemoryCatalog) SetHardwareInfo(agentID int, hwInfo *hwinfo.HWInfo) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if hwInfo == nil {
		delete(mc.hwInfoMap, agentID)
		return
	}
	mc.hwInfoMap[agentID] = hwInfo
}
