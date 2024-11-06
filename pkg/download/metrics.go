package download

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dockerPullSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_pull_success",
			Help: "Number of successful Docker image pulls.",
		},
		[]string{"image"},
	)

	dockerPullFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_pull_failure",
			Help: "Number of failed Docker image pulls.",
		},
		[]string{"image", "reason"},
	)

	dockerPullDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "docker_pull_duration_seconds",
			Help:    "Histogram of Docker image pull durations in seconds.",
			Buckets: []float64{1, 5, 10, 30, 60, 120}, // Adjust buckets as needed
		},
		[]string{"image"},
	)

	dockerPullSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "docker_pull_speed_bytes_per_second",
			Help: "Speed of Docker image pull in bytes per second.",
		},
		[]string{"image"},
	)

	dockerImageDeletionSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_image_deletion_success",
			Help: "Number of successful Docker image deletions.",
		},
		[]string{"image"},
	)

	dockerImageDeletionFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docker_image_deletion_failure",
			Help: "Number of failed Docker image deletions.",
		},
		[]string{"image", "reason"},
	)
	// imageToPull = "YOUR_DOCKER_IMAGE_HERE" // Update this (e.g., "nginx:latest")
)

func init() {
	prometheus.MustRegister(
		dockerPullSuccess,
		dockerPullFailure,
		dockerPullDuration,
		dockerPullSpeed,
		dockerImageDeletionSuccess,
		dockerImageDeletionFailure,
	)
}

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

	// download image
	PullDockerImage(context.Background(), "nginx:latest")
	PullDockerImage(context.Background(), "alpine:latest")
	PullDockerImage(context.Background(), "busybox:latest")
}
