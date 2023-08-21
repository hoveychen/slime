/*
Copyright © 2023 Harry C <hoveychen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package hwinfo

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/jaypipes/ghw"
)

type HWInfo struct {
	CPUCores         int
	CPUThreads       int
	MemoryUsableGB   float32
	MemoryPhysicalGB float32
	GPUNames         []string
	MacAddresses     []string
	PlatformArch     string
	PlatformOS       string
}

func NewHWInfo() *HWInfo {
	info := &HWInfo{
		PlatformArch: runtime.GOARCH,
		PlatformOS:   runtime.GOOS,
	}

	safeExec(func() {
		mem, err := ghw.Memory(ghw.WithDisableWarnings())
		if err != nil {
			return
		}

		info.MemoryUsableGB = float32(mem.TotalUsableBytes) / 1024 / 1024 / 1024
		info.MemoryPhysicalGB = float32(mem.TotalPhysicalBytes) / 1024 / 1024 / 1024
	})

	safeExec(func() {
		cpu, err := ghw.CPU(ghw.WithDisableWarnings())
		if err != nil {
			return
		}
		info.CPUCores = int(cpu.TotalCores)
		info.CPUThreads = int(cpu.TotalThreads)
	})

	safeExec(func() {
		gpu, err := ghw.GPU(ghw.WithDisableWarnings())
		if err != nil {
			return
		}
		for _, card := range gpu.GraphicsCards {
			if card.DeviceInfo != nil && card.DeviceInfo.Product != nil {
				info.GPUNames = append(info.GPUNames, card.DeviceInfo.Product.Name)
			} else {
				info.GPUNames = append(info.GPUNames, "Unknown")
			}
		}
		// compress the same GPU
		names := make(map[string]int)
		for _, name := range info.GPUNames {
			names[name]++
		}
		var newNames []string
		for k, v := range names {
			newName := k
			if v > 1 {
				newName = fmt.Sprintf("%dx %s", v, newName)
			}
			newNames = append(newNames, newName)
		}
		sort.Strings(newNames)
		info.GPUNames = newNames
	})

	safeExec(func() {
		network, err := ghw.Network(ghw.WithDisableWarnings())
		if err != nil {
			return
		}
		for _, nic := range network.NICs {
			if nic.IsVirtual {
				continue
			}

			info.MacAddresses = append(info.MacAddresses, nic.MacAddress)
		}
	})

	return info
}

func safeExec(f func()) {
	defer func() { recover() }()
	f()
}
