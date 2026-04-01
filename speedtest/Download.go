package speedtest

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func RunDownload() float64 {

	download, stopdownload := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopdownload()

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
	// run sequentially --same isp pipeline , tcp acks of download use upload pipeline , simple/clear
	// Workers running parallely to avoid RTT bottleneck
	numworkers := 4
	for i := 1; i <= numworkers; i++ {
		wg.Add(1)
		go GetDownload(download, &client, &totalBytes, &wg, i)

	}
	wg.Wait()

	downloadedBytes := totalBytes.Load()
	megabits := float64(downloadedBytes) * 8 / 1e6
	downloadSpeed := megabits / 30.0
	return downloadSpeed

}
func GetDownload(ctx context.Context, client *http.Client, totalBytes *atomic.Uint64, wg *sync.WaitGroup, wid int) {
	defer wg.Done()
	for {
		// Exit if Timeout
		if ctx.Err() != nil {
			return
		}
		// Create a get request to download 99 mb(cloudfare limit)
		req, err := http.NewRequestWithContext(ctx, "GET", "https://speed.cloudflare.com/__down?bytes=99999999", nil)
		if err != nil {
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Println("Worker error", wid, err)
			continue
		}
		// Track the number of bytes read and add it to the atomic var
		n, _ := io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		totalBytes.Add(uint64(n))

	}

}
