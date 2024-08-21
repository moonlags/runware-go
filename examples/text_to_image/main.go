package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/moonlags/runware-go"
)

func main() {
	client, err := runware.New(context.Background(), os.Getenv("RUNWARE_KEY"))
	if err != nil {
		log.Fatal("Can not create runware client:", err)
	}
	defer client.Close()

	images, err := client.TextToImage(runware.TextToImageArgs{
		PositivePrompt: "cool cat in sunglasses",
		Model:          "runware:100@1",
		IncludeCost:    true,
		NumberResults:  4,
	})
	if err != nil {
		log.Fatal("Can not generate image:", err)
	}

	for i, image := range images {
		fmt.Printf("%d: %s - %f\n", i, image.URL, image.Cost)
	}

	images, err = client.TextToImage(runware.TextToImageArgs{
		PositivePrompt: "a golden tree",
		Model:          "runware:100@1",
	})
	if err != nil {
		log.Fatal("Can not generate image:", err)
	}

	for i, image := range images {
		fmt.Printf("%d: %s - %f\n", i, image.URL, image.Cost)
	}
}
