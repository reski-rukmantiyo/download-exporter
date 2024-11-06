package download

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	fileDownloadSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_download_success",
			Help: "Number of successful file downloads.",
		},
		[]string{"filename"},
	)

	fileDownloadFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_download_failure",
			Help: "Number of failed file downloads.",
		},
		[]string{"filename", "reason"},
	)

	fileHashMismatch = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_hash_mismatch",
			Help: "Number of times the downloaded file's hash didn't match the expected hash.",
		},
		[]string{"filename", "expected_hash", "actual_hash"},
	)

	// expectedFileHash = "YOUR_EXPECTED_SHA256_HASH_HERE" // Update this
	// downloadURL      = "YOUR_FILE_DOWNLOAD_URL_HERE"    // Update this
	// filename         = "example.txt"                    // Update if needed
)

func Download() {
	log.Print("Downloading file")
}
