package monitoring

import (
	"log"
	"runtime"
	"time"
)

func MonitorSystemMetrics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
		log.Printf("TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
		log.Printf("Sys = %v MiB", bToMb(m.Sys))
		log.Printf("NumGC = %v", m.NumGC)
		log.Printf("Goroutines = %v", runtime.NumGoroutine())
		numGoroutine := runtime.NumGoroutine()
		log.Printf("Number of Goroutines: %d", numGoroutine)

		time.Sleep(1 * time.Minute)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
