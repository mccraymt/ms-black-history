package diagnostics

import "runtime"

// memStat is just a smaller version of GoLang runtime.MemStats, edited for easier reading.
type memStat struct {
	Alloc         uint64
	TotalAlloc    uint64
	Sys           uint64
	Lookups       uint64
	Mallocs       uint64
	Frees         uint64
	HeapAlloc     uint64
	HeapSys       uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapReleased  uint64
	HeapObjects   uint64
	StackInuse    uint64
	StackSys      uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	BuckHashSys   uint64
	GCSys         uint64
	OtherSys      uint64
	NextGC        uint64
	LastGC        uint64
	PauseTotalNs  uint64
	NumGC         uint32
	NumForcedGC   uint32
	GCCPUFraction float64
	EnableGC      bool
	DebugGC       bool
}

func abbreviate(src *runtime.MemStats) *memStat {
	if src == nil {
		return nil
	}

	rtn := memStat{
		Alloc:         src.Alloc,
		TotalAlloc:    src.TotalAlloc,
		Sys:           src.Sys,
		Lookups:       src.Lookups,
		Mallocs:       src.Mallocs,
		Frees:         src.Frees,
		HeapAlloc:     src.HeapAlloc,
		HeapSys:       src.HeapSys,
		HeapIdle:      src.HeapIdle,
		HeapInuse:     src.HeapInuse,
		HeapReleased:  src.HeapReleased,
		HeapObjects:   src.HeapObjects,
		StackInuse:    src.StackInuse,
		StackSys:      src.StackSys,
		MSpanInuse:    src.MSpanInuse,
		MSpanSys:      src.MSpanSys,
		MCacheInuse:   src.MCacheInuse,
		MCacheSys:     src.MCacheSys,
		BuckHashSys:   src.BuckHashSys,
		GCSys:         src.GCSys,
		OtherSys:      src.OtherSys,
		NextGC:        src.NextGC,
		LastGC:        src.LastGC,
		PauseTotalNs:  src.PauseTotalNs,
		NumGC:         src.NumGC,
		NumForcedGC:   src.NumForcedGC,
		GCCPUFraction: src.GCCPUFraction,
		EnableGC:      src.EnableGC,
		DebugGC:       src.DebugGC,
	}

	return &rtn
}
