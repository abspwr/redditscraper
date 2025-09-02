package redditscraper

import (
	"github.com/mmcdole/gofeed"
)

type Topic struct {
	Title string
	Url   string
}

func GetHotTopics(url string) (topics []Topic, err error) {
	fp := gofeed.NewParser()

	feeds, _ := fp.ParseURL(url)

	for _, f := range feeds.Items {
		topics = append(topics, Topic{Title: f.Title, Url: f.Link})
	}

	return topics, nil
}
