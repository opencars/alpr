package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/logger"
	"github.com/opencars/alpr/pkg/recognizer/openalpr"
)

var green = color.RGBA{
	R: 0, G: 255, B: 0, A: 255,
}

func main() {
	var configPath, imageURL string

	flag.StringVar(&configPath, "config", "./config/config.toml", "Path to the configuration file")
	flag.StringVar(&imageURL, "url", "", "URL to the image")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatalf("failed to read config: %v", err)
	}

	recognizer, err := openalpr.New(&conf.OpenALPR)
	if err != nil {
		logger.Fatalf("failed to initialize recognizer: %v", err)
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	var body bytes.Buffer
	dup := io.TeeReader(resp.Body, &body)

	results, err := recognizer.Recognize(dup)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	img, _, err := image.Decode(&body)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	m := image.NewRGBA(img.Bounds())
	draw.Draw(m, img.Bounds(), img, image.Point{}, draw.Src)

	for _, res := range results {
		prev := res.Coordinates[len(res.Coordinates)-1]

		for i, point := range res.Coordinates {
			for ; prev.X > point.X; prev.X-- {
				rect := image.Rect(prev.X, prev.Y, prev.X+3, prev.Y+3)
				draw.Draw(m, rect.Bounds(), &image.Uniform{C: green}, image.Point{X: prev.X, Y: prev.Y}, draw.Src)
			}

			for ; prev.X < point.X; prev.X++ {
				rect := image.Rect(prev.X, prev.Y, prev.X+3, prev.Y+3)
				draw.Draw(m, rect.Bounds(), &image.Uniform{C: green}, image.Point{X: prev.X, Y: prev.Y}, draw.Src)
			}

			for ; prev.Y > point.Y; prev.Y-- {
				rect := image.Rect(prev.X, prev.Y, prev.X+3, prev.Y+3)
				draw.Draw(m, rect.Bounds(), &image.Uniform{C: green}, image.Point{X: prev.X, Y: prev.Y}, draw.Src)
			}

			for ; prev.Y < point.Y; prev.Y++ {
				rect := image.Rect(prev.X, prev.Y, prev.X+3, prev.Y+3)
				draw.Draw(m, rect.Bounds(), &image.Uniform{C: green}, image.Point{X: prev.X, Y: prev.Y}, draw.Src)
			}

			prev = res.Coordinates[i]
		}
	}

	toimg, _ := os.Create("new.jpg")
	defer toimg.Close()

	err = jpeg.Encode(toimg, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

}
