package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	// // Create a 32kb slice of random data to act as our upload payload.
	payloadSize := 25 * 1024 * 1024
	payload := make([]byte, payloadSize)
	rand.Read(payload) // Fill the slice with random bytes

	download, stopdownload := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopdownload()

	// atomic variable to track total bytes
	var totalBytes atomic.Uint64
	var totalUploaded atomic.Uint64
	var wg sync.WaitGroup

	// Custom client to maintain Open connections
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}
	// run sequentially --same isp pipeline , tcp acks of download use upload pipeline , simple/clear
	// Workers running parallely to avoid RTT bottleneck
	numworkers := 1
	for i := 1; i <= numworkers; i++ {
		wg.Add(1)
		go GetDownload(download, &client, &totalBytes, &wg, i)

	}
	wg.Wait()

	log.Println("Completed Download Testing Starting Upload Speed Testing")

	downloadedBytes := totalBytes.Load()
	megabits := float64(downloadedBytes) * 8 / 1e6
	downloadSpeed := megabits / 30.0
	log.Printf("Download speed: %.2f Mbps\n", downloadSpeed)

	upload, stopupload := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopupload()

	for i := 1; i <= numworkers; i++ {
		wg.Add(1)
		go GetUpload(upload, &client, payload, &totalUploaded, &wg, i)
	}

	wg.Wait()

	// Calculating Download Speed

	uploadedBytes := totalUploaded.Load()
	meg := float64(uploadedBytes) * 8 / 1e6
	uploadedSpeed := meg / 30.0

	log.Printf("Upload Speed: %.2f", uploadedSpeed)
}
