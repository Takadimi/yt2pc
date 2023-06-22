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
	MimeType      string `json:"mimeType"`
	URL           string `json:"url"`
	ContentLength string `json:"contentLength"`
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
	Context inntertubeContext `json:"context"`
	VideoID string            `json:"videoId"`
}

type inntertubeContext struct {
	Client innertubeClient `json:"client"`
}

type innertubeClient struct {
	BrowserName    string `json:"browserName"`
	BrowserVersion string `json:"browserVersion"`
	ClientName     string `json:"clientName"`
	ClientVersion  string `json:"clientVersion"`
}

func VideoDataByInnertube(c *http.Client, id string) ([]byte, error) {
	// seems like same token for all WEB clients
	const webToken = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
	u := fmt.Sprintf("https://www.youtube.com/youtubei/v1/player?key=%s", webToken)

	data := innertubeRequest{
		Context: inntertubeContext{
			Client: innertubeClient{
				BrowserName:    "Mozilla",
				BrowserVersion: "5.0",
				ClientName:     "WEB",
				ClientVersion:  "2.20210617.01.00",
			},
		},
		VideoID: id,
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
