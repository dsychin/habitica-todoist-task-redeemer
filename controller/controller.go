package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dsychin/habitica-todoist-task-redeemer/config"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type (
	TodoistEvent struct {
		EventName EventName `json:"event_name"`
		UserID    string    `json:"user_id"`
		EventData EventData `json:"event_data"`
	}

	EventName string

	EventData struct {
		ID      string `json:"id"`
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

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		return events.APIGatewayProxyResponse{}, errors.New("empty payload")
	}

	var event TodoistEvent
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.Wrapf(err, "error unmarshalling json: %s", []byte(request.Body))
	}

	log.Printf("Received event %s", event.EventName)

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
		"notes": fmt.Sprintf("https://todoist.com/showTask?id=%s", event.EventData.ID),
	}

	log.Printf("Processing ID %s - text %s", event.EventData.ID, event.EventData.Content)

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
	req.Header.Set("x-api-user", config.UserID)
	req.Header.Set("x-api-key", config.APIToken)

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
	req.Header.Set("x-api-user", config.UserID)
	req.Header.Set("x-api-key", config.APIToken)

	res, err = client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Error in score task http POST")
	}
	defer res.Body.Close()

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "Error reading score task response body")

	}
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status is not 200 OK. Body: %s Request: %s", resBody, taskJson))
	}

	log.Printf("Status %d", res.StatusCode)

	return nil
}
