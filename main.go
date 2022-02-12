package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	HabiticaBaseURL string = "https://habitica.com/api/v3/"
	// AppID is the unique app id for this app
	AppID           string    = "b0fbabc9-32be-49c7-8b3e-f05578f61388-Todoist-Habitica-Task-Redeemer"
	ItemCompleted   EventName = "item:completed"
	ItemUncompleted EventName = "item:uncompleted"
)

var (
	UserID   string
	APIToken string
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
			return events.APIGatewayProxyResponse{}, errors.Wrap(err, "error in handleItemCompleted")
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
	client := &http.Client{}

	// TODO: handle tags to identify which task were created by this program

	// create task
	task := map[string]interface{}{
		"text":  event.EventData.Content,
		"type":  "todo",
		"notes": fmt.Sprintf("Todoist: %d", event.EventData.ID),
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		return errors.Wrapf(err, "Error marshaling taskJson")
	}
	jsonReader := bytes.NewReader(taskJson)

	// submit task create request
	taskCreateUrl := HabiticaBaseURL + "tasks/user"
	req, err := http.NewRequest(http.MethodPost, taskCreateUrl, jsonReader)
	if err != nil {
		return errors.Wrapf(err, "Error in new task create request. Event: %+v", event)
	}

	req.Header.Add("x-client", AppID)
	req.Header.Add("x-api-user", UserID)
	req.Header.Add("x-api-key", APIToken)

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Error in task create http POST")
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New("Status is not 201")
	}

	// complete task

	return nil
}

func main() {
	flag.String("habitica_key", "", "Your Habitica API Token")
	flag.String("habitica_user_id", "", "Your Habitica User ID")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("redeem")
	viper.BindEnv("habitica_key")
	viper.BindEnv("habitica_user_id")

	APIToken = viper.GetString("habitica_key")
	UserID = viper.GetString("habitica_user_id")

	lambda.Start(handler)
}
