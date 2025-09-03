package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	. "redditscraper/internal"
	"strings"
	"time"
)

const linkBase = "https://www.reddit.com/"
const linkTop = "https://www.reddit.com/best/communities/1/"
const linkRss = "/hot/.rss"
const dataname string = "datar.txt"
const inputename string = "inr.txt"

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

// func main() {

// 	// Make HTTP request
// 	// res, err := http.Get(linkTop)
// 	//res, err := http.NewRequest("GET", linkTop, nil)
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", linkTop, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
// 	}

// 	// Load HTML document
// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	count := 0

// 	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
// 		href, exists := s.Attr("href")
// 		if exists {
// 			if strings.HasPrefix(href, "/r/") {
// 				fmt.Println(href)
// 				if topics, err := GetHotTopics(linkBase + href + linkRss); err == nil {
// 					for _, t := range topics {
// 						appendLine(dataname, linkBase+href, t.Url, t.Title)
// 					}
// 				}
// 				time.Sleep(time.Duration(minSleepSeconds) * time.Second)
// 				count++
// 				if count >= 10 {
// 					return false
// 				}
// 			}
// 		}
// 		return true
// 	})

// }

func loadFile(filePath string) (out []string, err error) {

	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)

		out = append(out, fields[0])
	}

	return out, nil
}

func main() {

	log.Println("reading links from file")
	links, err := loadFile(inputename)

	if err != nil {
		log.Fatal("can't open input file")
	}

	for _, l := range links {
		log.Println("getting trendings topics from", l)
		if topics, err := GetHotTopics(l + linkRss); err == nil {
			for _, t := range topics {
				log.Println("writing to file", inputename)
				appendLine(dataname, l+linkRss, t.Url, t.Title)
			}
			time.Sleep(time.Duration(minSleepSeconds) * time.Second)
		}
	}

}
