package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
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
		// req.ContentLength = int64(len(payload))
		// req.Header.Set("Content-Type", "application/octet-stream")

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
