package main

import (
	"context"
	"log"
	"strings"

	"github.com/Takadimi/yt2pc/core/youtube"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	videoID, hasVideoID := request.PathParameters["id"]
	videoID = strings.TrimSpace(videoID)
	if !hasVideoID || videoID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "missing or invalid required path parameter 'id'",
		}, nil
	}

	videoData, err := youtube.GetVideoData(ctx, videoID)
	if err != nil {
		return serverErrorResponse("failed to get video data", err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 307, // temporary redirect
		Headers: map[string]string{
			"Location": videoData.Audio.URL,
		},
	}, nil
}

func serverErrorResponse(description string, err error) events.APIGatewayProxyResponse {
	log.Printf("%s: %s", description, err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       description,
	}
}

func main() {
	lambda.Start(Handler)
}
