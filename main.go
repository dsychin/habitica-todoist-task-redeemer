package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type (
	TodoistEvent struct {
		EventName EventName `json:"event_name"`
		UserID    int       `json:"user_id"`
		EventData EventData `json:"event_data"`
	}

	EventName string

	EventData struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	}
)

const (
	ItemCompleted   EventName = "item:completed"
	ItemUncompleted EventName = "item:uncompleted"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		return events.APIGatewayProxyResponse{}, errors.New("empty payload")
	}

	var event TodoistEvent
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("json unmarshal error")
	}

	switch event.EventName {
	case ItemCompleted:
		err := handleItemCompleted(event)
		if err != nil {
			return events.APIGatewayProxyResponse{}, errors.New("error in handleItemCompleted")
		}
	default:
		return events.APIGatewayProxyResponse{}, errors.New("invalid event name")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello AWS Lambda and Netlify",
	}, nil
}

func handleItemCompleted(event TodoistEvent) error {
	// create task

	// complete task

	return nil
}

func main() {
	lambda.Start(handler)
}
