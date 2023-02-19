package controller

import (
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
}
