package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	. "redditscraper/internal"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const linkBase = "https://www.reddit.com/"
const linkTop = "https://www.reddit.com/best/communities/1/"
const linkRss = "/hot/.rss"
const dataname string = "datar.txt"

const minSleepSeconds int = 5

const maxSubs = 10

func appendLine(filePath string, messages ...string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, strings.Join(messages, " "))
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

func main() {

	// Make HTTP request
	res, err := http.Get(linkTop)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if exists {
			if strings.HasPrefix(href, "/r/") {
				fmt.Println(href)
				if topics, err := GetHotTopics(linkBase + href + linkRss); err == nil {
					for _, t := range topics {
						appendLine(dataname, linkBase+href, t.Url, t.Title)
					}
				}
				time.Sleep(time.Duration(minSleepSeconds) * time.Second)
				count++
				if count >= 10 {
					return false
				}
			}
		}
		return true
	})

}
