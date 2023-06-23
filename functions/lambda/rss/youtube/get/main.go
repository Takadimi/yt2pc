package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Takadimi/yt2pc/core/http"
	"github.com/Takadimi/yt2pc/core/rss"
	"github.com/Takadimi/yt2pc/core/youtube"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	videoID, hasVideoID := request.PathParameters["id"]
	videoID = strings.TrimSpace(videoID)
	if !hasVideoID || videoID == "" {
		return http.BadRequest("missing or invalid required path parameter 'id'"), nil
	}

	videoData, err := youtube.GetVideoData(ctx, videoID)
	if err != nil {
		return http.ServerError("failed to get video data", err), nil
	}

	audioStream := videoData.GetAudioStream([]youtube.AudioPreference{
		{MIMEType: "audio/mp4", AudioQuality: youtube.AudioQualityLow},
		{MIMEType: "audio/mp4", AudioQuality: youtube.AudioQualityMedium},
		{MIMEType: "audio/mp4", AudioQuality: youtube.AudioQualityHigh},
	})
	if audioStream.URL == "" {
		return http.ServerError("no matching audio stream for video", err), nil
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
		return http.ServerError("failed to create RSS feed", err), nil
	}

	if request.RequestContext.HTTP.Method == "HEAD" {
		return http.OK(
			map[string]string{
				"Content-Type": "text/xml",
			},
			"",
		), nil
	}

	return http.OK(
		map[string]string{
			"Content-Type": "text/xml",
		},
		rssFeedXML,
	), nil
}

func main() {
	lambda.Start(Handler)
}
