package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// Imports the Google Cloud Vision API client package.
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name of the image file to annotate.
	filename := "static/cat.png"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	//	labels, err := client.DetectLabels(ctx, image, nil, 10)
	//	if err != nil {
	//		log.Fatalf("Failed to detect labels: %v", err)
	//	}
	web, err := client.DetectWeb(ctx, image, nil)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}
	raw, _ := json.MarshalIndent(web, "", "  ")
	fmt.Print(string(raw))
}
