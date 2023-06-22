package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type innertubeResponse struct {
	StreamingData innertubeStreamingData `json:"streamingData"`
	VideoDetails  innertubeVideoDetails  `json:"videoDetails"`
}

type innertubeStreamingData struct {
	AdaptiveFormats []innertubeAdaptiveFormat `json:"adaptiveFormats"`
}

type innertubeAdaptiveFormat struct {
	MIMEType         string `json:"mimeType"`
	URL              string `json:"url"`
	ContentLength    string `json:"contentLength"`
	AudioQuality     string `json:"audioQuality"`
	Bitrate          int    `json:"bitrate"`
	ApproxDurationMs string `json:"approxDurationMs"`
}

type innertubeVideoDetails struct {
	Title            string             `json:"title"`
	ShortDescription string             `json:"shortDescription"`
	Author           string             `json:"author"`
	ChannelID        string             `json:"channelId"`
	Thumbnail        innertubeThumbnail `json:"thumbnail"`
}

type innertubeThumbnail struct {
	Thumbnails []innertubeThumbnailItem `json:"thumbnails"`
}

type innertubeThumbnailItem struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type innertubeRequest struct {
	VideoID         string                    `json:"videoId,omitempty"`
	BrowseID        string                    `json:"browseId,omitempty"`
	Continuation    string                    `json:"continuation,omitempty"`
	Context         inntertubeContext         `json:"context"`
	PlaybackContext *innertubePlaybackContext `json:"playbackContext,omitempty"`
	ContentCheckOK  bool                      `json:"contentCheckOk,omitempty"`
	RacyCheckOk     bool                      `json:"racyCheckOk,omitempty"`
	Params          string                    `json:"params"`
}

type inntertubeContext struct {
	Client innertubeClient `json:"client"`
}

type innertubePlaybackContext struct {
	ContentPlaybackContext innertubeContentPlaybackContext `json:"contentPlaybackContext"`
}

type innertubeContentPlaybackContext struct {
	HTML5Preference string `json:"html5Preference"`
}

type innertubeClientInfo struct {
	Key    string
	Client innertubeClient
}

type innertubeClient struct {
	HL                string `json:"hl"`
	GL                string `json:"gl"`
	ClientName        string `json:"clientName"`
	ClientVersion     string `json:"clientVersion"`
	AndroidSDKVersion int    `json:"androidSDKVersion,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
	TimeZone          string `json:"timeZone"`
	UTCOffset         int    `json:"utcOffsetMinutes"`
}

// Leverage Android's client key and data
var clientInfo = innertubeClientInfo{
	Key: "AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w",
	Client: innertubeClient{
		HL:                "en",
		GL:                "US",
		TimeZone:          "UTC",
		ClientName:        "ANDROID",
		ClientVersion:     "17.31.35",
		UserAgent:         "com.google.android.youtube/17.31.35 (Linux; U; Android 11) gzip",
		AndroidSDKVersion: 30,
	},
}

func VideoDataByInnertube(c *http.Client, id string) ([]byte, error) {
	u := fmt.Sprintf("https://www.youtube.com/youtubei/v1/player?key=%s", clientInfo.Key)

	data := innertubeRequest{
		Context: inntertubeContext{
			Client: clientInfo.Client,
		},
		VideoID:        id,
		ContentCheckOK: true,
		RacyCheckOk:    true,
		Params:         "8AEB",
		PlaybackContext: &innertubePlaybackContext{
			ContentPlaybackContext: innertubeContentPlaybackContext{
				HTML5Preference: "HTML5_PREF_WANTS",
			},
		},
	}

	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return io.ReadAll(resp.Body)
}
