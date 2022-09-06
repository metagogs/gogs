package admin

import "github.com/metagogs/gogs/component"

type systemMemory struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	Lookups      uint64 `json:"lookups"`
	Mallocs      uint64 `json:"mallocs"`
	Frees        uint64 `json:"frees"`
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects"`
	StackInuse   uint64 `json:"stack_inuse"`
	StackSys     uint64 `json:"stack_sys"`
}

type systemStatus struct {
	Running            bool                      `json:"running"`
	DebugMode          bool                      `json:"debug_mode"`
	SessionCount       int64                     `json:"session_count"`
	OnlineSessionCount int64                     `json:"online_session_count"`
	NumGoroutine       int                       `json:"num_goroutine"`
	NumCgoCall         int64                     `json:"num_cgo_call"`
	NumCPU             int                       `json:"num_cpu"`
	SystemLatency      int64                     `json:"system_latency"`
	SystemLatencyList  []int64                   `json:"system_latency_list"`
	UserLatency        map[int64]int64           `json:"user_latency"`
	Acceptors          map[string]string         `json:"acceptors"`
	EncodeType         string                    `json:"encode_type"`
	DecodeType         map[string]uint8          `json:"decode_type"`
	Components         []component.ComponentDesc `json:"components"`
	Memory             systemMemory              `json:"memory"`
	Env                []string                  `json:"env"`
}
