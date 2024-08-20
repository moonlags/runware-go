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

	url, err := client.TextToImage(runware.TextToImageArgs{
		PositivePrompt: "cool cat in sunglasses",
		Model:          "runware:100@1",
	})
	if err != nil {
		log.Fatal("Can not generate image:", err)
	}

	fmt.Println(url)
}
