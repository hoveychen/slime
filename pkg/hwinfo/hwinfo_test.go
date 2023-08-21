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
	"runtime"
	"testing"
)

func TestNewHWInfo(t *testing.T) {
	info := NewHWInfo()

	if info.PlatformArch != runtime.GOARCH {
		t.Errorf("Expected PlatformArch to be %s, but got %s", runtime.GOARCH, info.PlatformArch)
	}

	if info.PlatformOS != runtime.GOOS {
		t.Errorf("Expected PlatformOS to be %s, but got %s", runtime.GOOS, info.PlatformOS)
	}

	if runtime.GOOS != "darwin" {
		if info.CPUCores == 0 {
			t.Error("Expected CPUCores to be greater than 0")
		}

		if info.CPUThreads == 0 {
			t.Error("Expected CPUThreads to be greater than 0")
		}

		if info.MemoryUsableGB == 0 {
			t.Error("Expected MemoryUsableGB to be greater than 0")
		}

		if info.MemoryPhysicalGB == 0 {
			t.Error("Expected MemoryPhysicalGB to be greater than 0")
		}

		// if len(info.GPUNames) == 0 {
		// 	t.Error("Expected at least one GPU name")
		// }
	}
}
