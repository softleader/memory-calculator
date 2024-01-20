package calc

import "github.com/paketo-buildpacks/libjvm/helper"

type MemoryLimitPath struct {
	V1 string
	V2 string
}

func NewMemoryLimitPath() *MemoryLimitPath {
	return &MemoryLimitPath{
		V1: helper.DefaultMemoryLimitPathV1, // cgroup v1 的記憶體上限路徑
		V2: helper.DefaultMemoryLimitPathV2, // cgroup v2 的記憶體上限路徑
	}
}
