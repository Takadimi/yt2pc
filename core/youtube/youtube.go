package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Youtube struct {
	ddb *dynamodb.Client
}

func New() Youtube {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	return Youtube{
		ddb: svc,
	}
}

type Video struct {
	VideoID string
}

func (y *Youtube) CreateVideo(ctx context.Context, videoID string) error {
	item, err := attributevalue.MarshalMap(Video{VideoID: videoID})
	if err != nil {
		return fmt.Errorf("subscription marshal: %w", err)
	}
	pk := fmt.Sprintf("VIDEO#%s", videoID)
	sk := pk
	setKeys(pk, sk, item)

	output, err := y.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put item: %w", err)
	}

	if os.Getenv("DEBUG") == "1" {
		log.Println(output)
	}

	return nil
}

func (y *Youtube) GetAudioURLForVideo(ctx context.Context, videoID string) (string, error) {
	hc := http.DefaultClient

	// test youtube-dl video (10s) "BaW_jenozKc"
	videoData, err := videoDataByInnertube(hc, videoID)
	if err != nil {
		log.Fatal(err)
	}

	var resp innertubeResponse
	if err := json.Unmarshal(videoData, &resp); err != nil {
		log.Fatal(err)
	}

	var audioMP4URL string
	for _, format := range resp.StreamingData.AdaptiveFormats {
		if strings.Contains(format.MimeType, "audio/mp4") {
			audioMP4URL = format.Url
		}
	}

	if audioMP4URL == "" {
		return "", fmt.Errorf("no audio URL for video ID %q", videoID)
	}

	return audioMP4URL, nil
}

type Subscription struct {
	ChannelID string
	Filter    string
}

func (y *Youtube) CreateSubscription(ctx context.Context, sub Subscription) error {
	item, err := attributevalue.MarshalMap(sub)
	if err != nil {
		return fmt.Errorf("subscription marshal: %w", err)
	}

	output, err := y.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put item: %w", err)
	}

	if os.Getenv("DEBUG") == "1" {
		log.Println(output)
	}

	return nil
}

func (y *Youtube) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	paginator := dynamodb.NewScanPaginator(y.ddb, &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	})

	items := []Subscription{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list subscriptions pagination: %w", err)
		}

		var pItems []Subscription
		err = attributevalue.UnmarshalListOfMaps(output.Items, &pItems)
		if err != nil {
			return nil, fmt.Errorf("list subscriptions unmarshal: %w", err)
		}
		items = append(items, pItems...)
	}

	return items, nil
}

func setKeys(pk, sk string, item map[string]types.AttributeValue) {
	item["PK"] = &types.AttributeValueMemberS{
		Value: pk,
	}
	item["SK"] = &types.AttributeValueMemberS{
		Value: sk,
	}
}

type innertubeResponse struct {
	StreamingData innertubeStreamingData `json:"streamingData"`
}

type innertubeStreamingData struct {
	AdaptiveFormats []innertubeAdaptiveFormat `json:"adaptiveFormats"`
}

type innertubeAdaptiveFormat struct {
	MimeType     string `json:"mimeType"`
	AudioQuality string `json:"audioQuality"`
	Url          string `json:"url"`
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

func videoDataByInnertube(c *http.Client, id string) ([]byte, error) {
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
