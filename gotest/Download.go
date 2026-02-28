package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

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
		// Track the number of bytes read and add it atomically
		n, _ := io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		totalBytes.Add(uint64(n))

	}

}
