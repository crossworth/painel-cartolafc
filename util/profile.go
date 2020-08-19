package util

import (
	"log"
	"time"
)

func Profile(name string, run func()) {
	start := time.Now()
	run()
	elapsed := time.Since(start)

	log.Printf("Profile %s: %s", name, elapsed)
}
