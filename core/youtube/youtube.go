package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type VideoData struct {
	ID           string
	URL          string
	Title        string
	Description  string
	Author       string
	ThumbnailURL string
	Audio        Audio
	AudioStreams []AudioStream
}

type Audio struct {
	URL           string
	MIMEType      string
	ContentLength string // size in bytes
}

type AudioStream struct {
	URL           string
	MIMEType      string
	FullMIMEType  string
	ContentLength string // size in bytes
	AudioQuality  AudioQuality
	Bitrate       int
	Duration      time.Duration
}

type AudioQuality int

const (
	AudioQualityUnknown AudioQuality = iota
	AudioQualityLow
	AudioQualityMedium
	AudioQualityHigh
)

type AudioPreference struct {
	MIMEType     string
	AudioQuality AudioQuality
}

func (videoData VideoData) GetAudioStream(preferences []AudioPreference) AudioStream {
	streamsByQuality := map[AudioPreference]AudioStream{}

	for _, stream := range videoData.AudioStreams {
		streamsByQuality[AudioPreference{
			MIMEType:     stream.MIMEType,
			AudioQuality: stream.AudioQuality,
		}] = stream
	}

	for _, preference := range preferences {
		if stream, hasStream := streamsByQuality[preference]; hasStream {
			return stream
		}
	}

	return AudioStream{}
}

var ErrVideoNotAvailable = errors.New("video not available")

func GetVideoData(ctx context.Context, videoID string) (VideoData, error) {
	innertubeVideoData, err := VideoDataByInnertube(http.DefaultClient, videoID)
	if err != nil {
		return VideoData{}, fmt.Errorf("GetVideoData innertube api call: %w", err)
	}

	var resp innertubeResponse
	if err := json.Unmarshal(innertubeVideoData, &resp); err != nil {
		return VideoData{}, fmt.Errorf("GetVideoData innertube response unmarshaling: %w", err)
	}

	if resp.PlayabilityStatus.Status != "OK" {
		return VideoData{}, ErrVideoNotAvailable
	}

	audioStreams := make([]AudioStream, 0)
	for _, format := range resp.StreamingData.AdaptiveFormats {
		if strings.HasPrefix(format.MIMEType, "audio") {
			var quality AudioQuality
			switch format.AudioQuality {
			case "AUDIO_QUALITY_LOW":
				quality = AudioQualityLow
			case "AUDIO_QUALITY_MEDIUM":
				quality = AudioQualityMedium
			case "AUDIO_QUALITY_HIGH":
				quality = AudioQualityHigh
			}

			durationInt, err := strconv.Atoi(format.ApproxDurationMs)
			if err != nil {
				return VideoData{}, fmt.Errorf("GetVideoData innertube response duration parsing: %w", err)
			}
			duration := time.Millisecond * time.Duration(durationInt)

			audioStreams = append(audioStreams, AudioStream{
				URL:           format.URL,
				MIMEType:      basicMIMEType(format.MIMEType),
				FullMIMEType:  format.MIMEType,
				ContentLength: format.ContentLength,
				AudioQuality:  quality,
				Bitrate:       format.Bitrate,
				Duration:      duration,
			})
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
		AudioStreams: audioStreams,
	}, nil
}

func url(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
}

func basicMIMEType(mimeType string) string {
	parts := strings.Split(mimeType, ";")
	return parts[0]
}
