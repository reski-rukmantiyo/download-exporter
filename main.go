package main

import (
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

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}
