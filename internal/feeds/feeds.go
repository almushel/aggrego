package feeds

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

type ChannelItem struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	PubDate       string `xml:"pubDate,omitempty"`
	Guid          string `xml:"guid,omitempty"`
	Description   string `xml:"description"`
	Generator     string `xml:"generator,omitempty"`
	Language      string `xml:"language,omitempty"`
	LastBuildDate string `xml:"lastBuildDate,omitempty"`
	AtomLink      string `xml:"atom:link,omitempty"`
}

type RSSChannel struct {
	ChannelItem
	Items []ChannelItem `xml:"item"`
}

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

func FetchRSSFeed(url string) (RSSFeed, error) {
	var result RSSFeed

	response, err := http.Get(url)
	if err != nil {
		return result, err
	} else if response.StatusCode != 200 {
		return result, errors.New(response.Status)
	}

	buf, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return result, err
	}

	err = xml.Unmarshal(buf, &result)
	return result, err
}
