package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Takadimi/yt2pc/core/rss"
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

	audioURL := fmt.Sprintf("https://yq0hieyjqa.execute-api.us-east-1.amazonaws.com/youtube/%s", videoID)
	rssFeedXML := rss.CreateSingleVideoPodcastFeed(videoID, audioURL)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/xml",
		},
		Body: rssFeedXML,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
