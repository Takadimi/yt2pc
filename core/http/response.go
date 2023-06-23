package http

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func OK(headers map[string]string, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
	}
}

func TemporaryRedirect(location string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 307,
		Headers: map[string]string{
			"Location": location,
		},
	}
}

func BadRequest(description string) events.APIGatewayProxyResponse {
	log.Println(description)
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       description,
	}
}

func ServerError(description string, err error) events.APIGatewayProxyResponse {
	log.Printf("%s: %s", description, err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       description,
	}
}
