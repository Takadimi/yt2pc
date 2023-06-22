package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Takadimi/yt2pc/core/rss"
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

	audioStream := videoData.GetAudioStream([]youtube.AudioPreference{
		{"audio/mp4", youtube.AudioQualityLow},
		{"audio/mp4", youtube.AudioQualityMedium},
		{"audio/mp4", youtube.AudioQualityHigh},
	})
	if audioStream.URL == "" {
		return serverErrorResponse("no matching audio stream for video", err), nil
	}

	proxyAudioURL := fmt.Sprintf("%s/youtube/%s", os.Getenv("API_ENDPOINT"), videoID)
	rssFeedXML, err := rss.CreateSingleVideoPodcastFeed(
		videoID,
		videoData.URL,
		videoData.Title,
		videoData.Description,
		videoData.Author,
		videoData.ThumbnailURL,
		proxyAudioURL,
		audioStream.MIMEType,
		audioStream.ContentLength,
		audioStream.Duration,
	)
	if err != nil {
		return serverErrorResponse("failed to create RSS feed", err), nil
	}

	if request.RequestContext.HTTP.Method == "HEAD" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text/xml",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/xml",
		},
		Body: rssFeedXML,
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
