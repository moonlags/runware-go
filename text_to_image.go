package runware

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type textToImageRequest struct {
	TaskType       string  `json:"taskType"`
	TaskUUID       string  `json:"taskUUID"`
	PositivePrompt string  `json:"positivePrompt"`
	NegativePrompt string  `json:"negativePrompt,omitempty"`
	OutputFormat   string  `json:"outputFormat,omitempty"`
	Model          string  `json:"model"`
	NumberResults  int     `json:"numberResults"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	Steps          int     `json:"steps,omitempty"`
	CFGScale       float32 `json:"CFGScale,omitempty"`
	IncludeCost    bool    `json:"includeCost"`
}

type TextToImageArgs struct {
	PositivePrompt string
	NegativePrompt string
	OutputFormat   string
	Model          string
	NumberResults  int
	Width          int
	Height         int
	Steps          int
	CFGScale       float32
	IncludeCost    bool
}

func (c *Client) TextToImage(args TextToImageArgs) ([]Image, error) {
	req, err := ttIArgsToRequest(args)
	if err != nil {
		return nil, err
	}

	sendData, err := json.Marshal([]*textToImageRequest{req})
	if err != nil {
		return nil, nil
	}

	if err := c.Send(sendData); err != nil {
		return nil, nil
	}

	images := make([]Image, 0)

	for msg := range c.incomingMessages {
		var msgData socketMessage
		if err := json.Unmarshal(msg, &msgData); err != nil {
			return nil, err
		}

		if err := c.checkError(msgData); err != nil {
			return nil, err
		}

		if msgData.Data[0]["taskType"].(string) != "imageInference" || msgData.Data[0]["taskUUID"] != req.TaskUUID {
			c.incomingMessages <- msg
			continue
		}

		var cost float64
		if _, ok := msgData.Data[0]["cost"]; ok {
			cost = msgData.Data[0]["cost"].(float64)
		}

		images = append(images, Image{
			URL:  msgData.Data[0]["imageURL"].(string),
			Cost: cost,
		})

		if len(images) == req.NumberResults {
			break
		}
	}
	return images, nil
}

func ttIArgsToRequest(args TextToImageArgs) (*textToImageRequest, error) {
	if args.PositivePrompt == "" {
		return nil, fmt.Errorf("PositivePrompt is required")
	}

	if args.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	if args.NumberResults == 0 {
		args.NumberResults = 1
	}

	if args.Width == 0 {
		args.Width = 512
	}

	if args.Height == 0 {
		args.Height = 512
	}

	return &textToImageRequest{
		TaskType:       "imageInference",
		TaskUUID:       uuid.NewString(),
		PositivePrompt: args.PositivePrompt,
		NegativePrompt: args.NegativePrompt,
		OutputFormat:   args.OutputFormat,
		Model:          args.Model,
		NumberResults:  args.NumberResults,
		Width:          args.Width,
		Height:         args.Height,
		Steps:          args.Steps,
		CFGScale:       args.CFGScale,
		IncludeCost:    args.IncludeCost,
	}, nil
}
