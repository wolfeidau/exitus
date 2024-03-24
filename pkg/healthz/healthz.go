package healthz

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// contains statistics about the Go's process.
type runtimeStats struct {
	CollectedAt      time.Time `json:"-"`
	Arch             string    `json:"arch"`
	OS               string    `json:"os"`
	Version          string    `json:"version"`
	GoroutinesCount  int       `json:"goroutines_count"`
	HeapObjectsCount int       `json:"heap_objects_count"`
	AllocBytes       int       `json:"alloc_bytes"`
	TotalAllocBytes  int       `json:"total_alloc_bytes"`
}

type healthCheckResult struct {
	Message string        `json:"message,omitempty"`
	Runtime *runtimeStats `json:"runtime,omitempty"`
}

// Handler simple health handler.
func Handler() http.HandlerFunc {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			log.Info().Object("rs", collect()).Msg("stats")
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		result := healthCheckResult{
			Message: "ok",
			Runtime: collect(),
		}

		data, _ := json.Marshal(result)
		_, _ = w.Write(data) // return 200 by default
	}
}

func collect() *runtimeStats {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	return &runtimeStats{
		CollectedAt:      time.Now(),
		Arch:             runtime.GOARCH,
		OS:               runtime.GOOS,
		Version:          runtime.Version(),
		GoroutinesCount:  runtime.NumGoroutine(),
		HeapObjectsCount: int(ms.HeapObjects),
		AllocBytes:       int(ms.Alloc),
		TotalAllocBytes:  int(ms.TotalAlloc),
	}
}

func (rs *runtimeStats) MarshalZerologObject(e *zerolog.Event) {
	e.Str("arch", rs.Arch).Str("os", rs.OS).Str("version", rs.Version).
		Int("goroutines_count", rs.GoroutinesCount).
		Int("heap_objects_count", rs.HeapObjectsCount).
		Int("alloc_bytes", rs.AllocBytes).
		Int("total_alloc_bytes", rs.TotalAllocBytes)
}
