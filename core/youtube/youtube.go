package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type VideoData struct {
	ID           string
	URL          string
	Title        string
	Description  string
	Author       string
	ThumbnailURL string
	Audio        Audio
}

type Audio struct {
	URL           string
	MIMEType      string
	ContentLength string // size in bytes
}

func GetVideoData(ctx context.Context, videoID string) (VideoData, error) {
	innertubeVideoData, err := VideoDataByInnertube(http.DefaultClient, videoID)
	if err != nil {
		return VideoData{}, fmt.Errorf("GetVideoData innertube api call: %w", err)
	}

	var resp innertubeResponse
	if err := json.Unmarshal(innertubeVideoData, &resp); err != nil {
		return VideoData{}, fmt.Errorf("GetVideoData innertube response unmarshaling: %w", err)
	}

	var audio Audio
	for _, format := range resp.StreamingData.AdaptiveFormats {
		if strings.Contains(format.MimeType, "audio/mp4") {
			audio.URL = format.URL
			audio.MIMEType = "audio/mp4"
			audio.ContentLength = format.ContentLength
			break
		}
	}

	var thumbnailURL string
	if len(resp.VideoDetails.Thumbnail.Thumbnails) > 0 {
		thumbnailURL = resp.VideoDetails.Thumbnail.Thumbnails[0].URL
	}

	return VideoData{
		ID:           videoID,
		URL:          url(videoID),
		Title:        resp.VideoDetails.Title,
		Description:  resp.VideoDetails.ShortDescription,
		Author:       resp.VideoDetails.Author,
		ThumbnailURL: thumbnailURL,
		Audio:        audio,
	}, nil
}

func url(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
}
