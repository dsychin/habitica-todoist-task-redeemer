package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	t.Run("should return error if request body is empty", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{}

		_, err := HandleRequest(req)

		assert.EqualError(t, err, "empty payload")
	})

	t.Run("should return error if request body is an invalid json", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{
			Body: "{{",
		}

		_, err := HandleRequest(req)

		assert.ErrorContains(t, err, "error unmarshalling json")
	})

	t.Run("should return error if there is an unexpected event name", func(t *testing.T) {
		e := TodoistEvent{
			EventName: "invalid_event",
		}
		body, _ := json.Marshal(e)
		req := events.APIGatewayProxyRequest{
			Body: string(body),
		}

		_, err := HandleRequest(req)

		assert.EqualError(t, err, "invalid event name")
	})

	t.Run("should return 200 OK when status is uncompleted", func(t *testing.T) {
		e := TodoistEvent{
			EventName: "item:uncompleted",
		}
		body, _ := json.Marshal(e)
		req := events.APIGatewayProxyRequest{
			Body: string(body),
		}

		res, err := HandleRequest(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
