package main

import (
	"encoding/json"
	"errors"
	"net/http"

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
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "empty payload",
		}, errors.New("empty payload")
	}

	var event TodoistEvent
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error",
		}, errors.New("json unmarshal error")
	}

	// TODO: create habitica task

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello AWS Lambda and Netlify",
	}, nil
}

func main() {
	lambda.Start(handler)
}
