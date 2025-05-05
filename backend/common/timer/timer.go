package timer

import (
	"log"
	"time"
)

func Track(start time.Time, name string, logger *log.Logger) {
	elapsed := time.Since(start)
	logger.Printf("%s took %s", name, elapsed)
}
