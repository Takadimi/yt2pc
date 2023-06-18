package main

import (
	"context"
	"log"
	"strings"

	"github.com/Takadimi/yt2pc/core/youtube"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var youtubeSvc youtube.Youtube

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	videoID, hasVideoID := request.PathParameters["id"]
	videoID = strings.TrimSpace(videoID)
	if !hasVideoID || videoID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "missing or invalid required path parameter 'id'",
		}, nil
	}

	audioURL, err := youtubeSvc.GetAudioURLForVideo(ctx, videoID)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "failed to get audio url",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 307, // temporary redirect
		Headers: map[string]string{
			"Location": audioURL,
		},
	}, nil
}

func main() {
	youtubeSvc = youtube.New()

	lambda.Start(Handler)
}
