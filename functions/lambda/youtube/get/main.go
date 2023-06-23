package main

import (
	"context"
	"strings"

	"github.com/Takadimi/yt2pc/core/http"
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

	if request.RequestContext.HTTP.Method == "HEAD" {
		return http.OK(
			map[string]string{
				"Content-Length": audioStream.ContentLength,
				"Content-Type":   audioStream.FullMIMEType,
			},
			"",
		), nil
	}

	return http.TemporaryRedirect(audioStream.URL), nil
}

func main() {
	lambda.Start(Handler)
}
