# Golang Runware SDK

The SDK is used to run image inference with the Runware API, powered by the RunWare inference platform. It can be used to generate imaged with text-to-image prompt

This is not an official library, but a simple and atleast working copy of it.

## Get Api Access

For an API Key and trial credits, Create a free account with Runware

## Usage
 `go get github.com/moonlags/runware-go`

### Basic text-to-image example
```Golang
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
```
For more examples see `examples` directory

This library uses websockets under the box, because runway does not support a http api yet
