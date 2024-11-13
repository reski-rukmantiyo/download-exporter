package main

import (
	"fmt"
	"log"
	"time"

	"github.com/reski-rukmantiyo/download-exporter/pkg/download"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

func main() {
	// initialise gofr object
	app := gofr.New()

	app.GET("/download", func(ctx *gofr.Context) (interface{}, error) {
		// get the file path from the query params
		download.Download()
		return response.Raw{Data: "Downloaded"}, nil
	})

	// register route greet
	app.GET("/", func(ctx *gofr.Context) (interface{}, error) {

		return "Hello Download Exporter!", nil
	})

	minToPull := download.ImageConfigData.MinuteToPull
	formatCronjob := fmt.Sprintf("*/%d * * * * *", minToPull)

	// Run the cron job every 1 hour
	app.AddCronJob(formatCronjob, "", func(ctx *gofr.Context) {
		log.Printf("Cron job running at %s", time.Now().String())

		if download.IsProcessing() {
			log.Printf("Download is in progress")
			return
		}
		download.Download()

		log.Printf("Cron job finished at %s", time.Now().String())
	})

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}
