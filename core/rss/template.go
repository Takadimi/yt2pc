package rss

import (
	"fmt"
	nurl "net/url"
	"strings"
	"text/template"
	"time"
)

const podcastTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:podcast="https://podcastindex.org/namespace/1.0">
  <channel>
    <title>{{.Title}}</title>
    <description>{{.Description}}</description>
    <podcast:guid>{{.ID}}</podcast:guid>
    <podcast:locked>yes</podcast:locked>
    <language>en</language>
    <pubDate>{{.PubDate}}</pubDate>
    <lastBuildDate>{{.BuildDate}}</lastBuildDate>
    <link>{{.Link}}</link>
    <image>
        <url>{{.ThumbnailURL}}</url>
        <title>{{.Title}}</title>
        <link>{{.Link}}</link>
    </image>
    <item>
      <title>{{.Title}}</title>
      <podcast:episode>1</podcast:episode>
      <guid isPermaLink="false">{{.ID}}</guid>
      <link>{{.AudioURL}}</link>
      <description>{{.Description}}</description>
      <pubDate>{{.PubDate}}</pubDate>
      <author>{{.Author}}</author>
      <enclosure url="{{.AudioURL}}" length="{{.AudioContentLength}}" type="{{.AudioMIMEType}}"/>
    </item>
  </channel>
</rss>`

type templateData struct {
	Title              string
	Description        string
	ID                 string
	PubDate            string
	BuildDate          string
	Link               string
	Author             string
	ThumbnailURL       string
	AudioURL           string
	AudioContentLength string
	AudioMIMEType      string
}

func CreateSingleVideoPodcastFeed(id, url, title, description, author, thumbnailURL, audioURL, audioMIMEType, audioContentLength string) (string, error) {
	t, err := template.New("podcastFeed").Parse(podcastTemplate)
	if err != nil {
		return "", fmt.Errorf("CreateSingleVideoPodcastFeed parse template: %w", err)
	}

	tu, err := nurl.Parse(thumbnailURL)
	if err != nil {
		return "", fmt.Errorf("CreateSingleVideoPodcastFeed parse thumbnail URL: %w", err)
	}
	tu.RawQuery = ""

	nowDateStr := time.Now().Format(time.RFC1123Z)
	td := templateData{
		Title:              title,
		Description:        description,
		ID:                 id,
		PubDate:            nowDateStr,
		BuildDate:          nowDateStr,
		Link:               url,
		Author:             author,
		ThumbnailURL:       tu.String(),
		AudioURL:           audioURL,
		AudioContentLength: audioContentLength,
		AudioMIMEType:      audioMIMEType,
	}

	var sb strings.Builder
	if err := t.Execute(&sb, td); err != nil {
		return "", fmt.Errorf("CreateSingleVideoPodcastFeed execute template: %w", err)
	}

	return sb.String(), nil
}
