package main

import (
	"log"

	"github.com/Abdu-Rauf/speedext/speedtest"
)

func main() {
	log.Println("Starting speedext test...")

	log.Println("--- Running Download Test ---")
	downloadResult := speedtest.RunDownload()
	log.Println(downloadResult)

	log.Println("--- Running Upload Test ---")
	uploadResult := speedtest.RunUpload()
	log.Println(uploadResult)

	log.Println("Test complete.")
}
