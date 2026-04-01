package speedtest

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// 1. The custom reader that counts bytes as the HTTP client pulls them
type trackingReader struct {
	r          io.Reader
	totalBytes *atomic.Uint64
}

// Read intercepts the data flow, adds the byte count to our total, and passes the data along
func (t *trackingReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		// Atomically add the bytes exactly as they are sent to the OS network buffer
		t.totalBytes.Add(uint64(n))
	}
	return n, err
}
func RunUpload() float64 {
	// // Create a 32kb slice of random data to act as our upload payload.
	payloadSize := 25 * 1024 * 1024
	payload := make([]byte, payloadSize)
	rand.Read(payload) // Fill the slice with random bytes

	// atomic variable to track total bytes
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
	numworkers := 4
	upload, stopupload := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopupload()

	for i := 1; i <= numworkers; i++ {
		wg.Add(1)
		go GetUpload(upload, &client, payload, &totalUploaded, &wg, i)
	}

	wg.Wait()

	// Calculating Upload Speed
	uploadedBytes := totalUploaded.Load()
	meg := float64(uploadedBytes) * 8 / 1e6
	uploadSpeed := meg / 30.0

	return uploadSpeed
}

// 2. The upload worker that continuously pushes data to Cloudflare
func GetUpload(ctx context.Context, client *http.Client, payload []byte, totalBytes *atomic.Uint64, wg *sync.WaitGroup, wid int) {
	defer wg.Done()

	for {
		if ctx.Err() != nil {
			return // The 30-second timer is up, exit the loop cleanly
		}

		// Create a fresh reader from our dummy payload for this specific HTTP request
		dataReader := bytes.NewReader(payload)

		// Wrap that data reader in our custom tracking reader
		tr := &trackingReader{
			r:          dataReader,
			totalBytes: totalBytes,
		}

		// Build the POST request using the tracking reader as the request body
		req, err := http.NewRequestWithContext(ctx, "POST", "https://speed.cloudflare.com/__up", tr)
		if err != nil {
			continue // If request building fails, try again
		}

		// Set the content type so Cloudflare knows we are just sending raw binary data(is this even necessary if so why??)
		req.ContentLength = int64(len(payload))
		req.Header.Set("Content-Type", "application/octet-stream")

		// Execute the upload
		resp, err := client.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return // Context expired mid-upload, exit cleanly
			}
			log.Printf("Worker %d network error: %v\n", wid, err)
			continue // Network hiccup, restart the loop and try uploading again
		}

		// Even on an upload, the server sends a tiny response body back.
		// We discard it and close it to prevent memory leaks.
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
