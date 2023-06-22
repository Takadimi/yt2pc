package rss

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

const podcastTemplateWithItunes = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:podcast="https://podcastindex.org/namespace/1.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>{{.Title}}</title>
    <description>{{.Description}}</description>
    <podcast:guid>{{.ID}}</podcast:guid>
    <podcast:locked>no</podcast:locked>
    <language>en</language>
    <pubDate>{{.PubDate}}</pubDate>
    <lastBuildDate>{{.BuildDate}}</lastBuildDate>
    <link>{{.Link}}</link>
    <image>
        <url>{{.ThumbnailURL}}</url>
        <title>{{.Title}}</title>
        <link>{{.Link}}</link>
    </image>
    <itunes:author>{{.Author}}</itunes:author>
    <itunes:subtitle>{{.Description}}</itunes:subtitle>
    <itunes:summary>{{.Description}}</itunes:summary>
    <itunes:image href="{{.ThumbnailURL}}" />
    <itunes:explicit>no</itunes:explicit>
    <itunes:category text="Technology" />
    <item>
      <title>{{.Title}}</title>
      <podcast:episode>1</podcast:episode>
      <guid isPermaLink="false">{{.ID}}</guid>
      <link>{{.AudioURL}}</link>
      <description>{{.Description}}</description>
      <pubDate>{{.PubDate}}</pubDate>
      <author>{{.Author}}</author>
      <itunes:author>{{.Author}}</itunes:author>
      <itunes:subtitle>{{.Description}}</itunes:subtitle>
      <itunes:summary>{{.Description}}</itunes:summary>
      <itunes:length>{{.AudioContentLength}}</itunes:length>
      <itunes:duration>{{.AudioDuration}}</itunes:duration>
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
	AudioDuration      string
}

func CreateSingleVideoPodcastFeed(id, url, title, description, author, thumbnailURL, audioURL, audioMIMEType, audioContentLength string, audioDuration time.Duration) (string, error) {
	t, err := template.New("podcastFeed").Parse(podcastTemplateWithItunes)
	if err != nil {
		return "", fmt.Errorf("CreateSingleVideoPodcastFeed parse template: %w", err)
	}

	nowDateStr := time.Now().Format(time.RFC1123Z)
	td := templateData{
		Title:              title,
		Description:        description,
		ID:                 id,
		PubDate:            nowDateStr,
		BuildDate:          nowDateStr,
		Link:               url,
		Author:             author,
		ThumbnailURL:       strings.Split(thumbnailURL, "?")[0],
		AudioURL:           audioURL,
		AudioContentLength: audioContentLength,
		AudioMIMEType:      audioMIMEType,
		AudioDuration:      hms(audioDuration),
	}

	var sb strings.Builder
	if err := t.Execute(&sb, td); err != nil {
		return "", fmt.Errorf("CreateSingleVideoPodcastFeed execute template: %w", err)
	}

	return sb.String(), nil
}

func hms(duration time.Duration) string {
	totalSecs := int(duration.Seconds())
	h := totalSecs / 3600
	totalSecs = totalSecs % 3600
	m := totalSecs / 60
	s := totalSecs % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
