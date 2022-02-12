package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mitchellh/mapstructure"
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

	HabiticaResponse struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}

	HabiticaTask struct {
		ID   string `json:"id"`
		Text string `json:"text"`
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
	case ItemUncompleted:
		// TODO: revert if item is uncompleted
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
	// TODO: split into smaller functions

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

	// submit task create request
	taskCreateUrl := HabiticaBaseURL + "tasks/user"
	req, err := http.NewRequest(http.MethodPost, taskCreateUrl, bytes.NewBuffer(taskJson))
	if err != nil {
		return errors.Wrapf(err, "Error in new task create request. Event: %+v", event)
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-client", AppID)
	req.Header.Set("x-api-user", UserID)
	req.Header.Set("x-api-key", APIToken)

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Error in task create http POST")
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "Error reading task create response body")

	}
	if res.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Status is not 201. Body: %s Request: %s", resBody, taskJson))
	}

	var habiticaRes HabiticaResponse
	err = json.Unmarshal(resBody, &habiticaRes)
	if err != nil {
		return errors.Wrapf(err, "Error unmarshaling into HabiticaResponse. Body: %s", resBody)
	}
	if !habiticaRes.Success {
		return errors.New(fmt.Sprintf("Score task response is not success. Response: %s", resBody))
	}

	var habiticaTask HabiticaTask
	err = mapstructure.Decode(habiticaRes.Data, &habiticaTask)
	if err != nil {
		return errors.Wrapf(err, "Error in converting Habitica response data to HabiticaTask. Response: %s", resBody)
	}

	// complete task
	scoreTaskUrl := HabiticaBaseURL + fmt.Sprintf("tasks/%s/score/up", habiticaTask.ID)
	req, err = http.NewRequest(http.MethodPost, scoreTaskUrl, nil)
	if err != nil {
		return errors.Wrapf(err, "Error in score task request. Event: %+v", event)
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-client", AppID)
	req.Header.Set("x-api-user", UserID)
	req.Header.Set("x-api-key", APIToken)

	res, err = client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Error in score task http POST")
	}
	defer res.Body.Close()

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "Error reading score task response body")

	}
	if res.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Status is not 201. Body: %s Request: %s", resBody, taskJson))
	}

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

	if APIToken == "" {
		log.Fatal("API Token is empty. Please set the REDEEM_HABITICA_KEY environment variable.")
	}
	if UserID == "" {
		log.Fatal("API Token is empty. Please set the REDEEM_HABITICA_USER_ID environment variable.")
	}

	lambda.Start(handler)
}
