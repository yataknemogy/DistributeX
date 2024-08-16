package metrics

import (
	"log"
	"time"
)

func CollectMetrics() {
	for {
		log.Println("Collecting metrics...")
		time.Sleep(10 * time.Second)
	}
}
