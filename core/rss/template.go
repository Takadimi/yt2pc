package rss

import (
	"fmt"
	"strings"
)

const podcastTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:podcast="https://podcastindex.org/namespace/1.0">
  <channel>
    <title>Test Title</title>
    <description>This is a test description.</description>
    <podcast:guid>778116ac-6b1e-5ae2-b037-26a7ff2aee64</podcast:guid>
    <podcast:locked owner="ethan.woodward@gmail.com">yes</podcast:locked>
    <language>en</language>
    <pubDate>Fri, 16 Jun 2023 03:35:09 -0700</pubDate>
    <lastBuildDate>Fri, 16 Jun 2023 03:48:51 -0700</lastBuildDate>
    <link>https://youtube.com</link>
    <item>
      <title>Test Podcast</title>
      <podcast:episode>1</podcast:episode>
      <guid isPermaLink="false">714980e9-a5b3-4ab6-b54b-aff455edd8c2</guid>
      <link>https://youtube.com/watch?v=%s</link>
      <description>
        <![CDATA[<p>This is only a <strong>test!</strong></p>]]>
      </description>
      <pubDate>Tue, 06 Jun 2023 17:55:16 -0700</pubDate>
      <author>Youtube</author>
      <enclosure url="%s" length="56062633" type="audio/mp4"/>
    </item>
  </channel>
</rss>
`

func CreateSingleVideoPodcastFeed(videoID, audioURL string) string {
	return strings.TrimSpace(fmt.Sprintf(podcastTemplate, videoID, audioURL))
}
