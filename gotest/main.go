package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// atomic variable to track total bytes
	var totalBytes atomic.Uint64
	var wg sync.WaitGroup

	// Custom client to maintain Open connections
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	// Workers running parallely to avoid RTT bottleneck
	numworkers := 4
	for i := 1; i <= numworkers; i++ {
		wg.Add(1)
		go GetDownload(ctx, &client, &totalBytes, &wg, i)

	}
	wg.Wait()

	// Calculating Download Speed
	downloadedBytes := totalBytes.Load()
	megabits := float64(downloadedBytes) * 8 / 1e6
	downloadSpeed := megabits / 30.0

	log.Printf("Test completed!\n Downloaded: %.2f Mbs\n Download speed: %.2f", megabits, downloadSpeed)
}
