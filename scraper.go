package rscraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//SubReddit contains urls of tabs of the given subreddit
type SubReddit struct {
	URL           string `json:"url"` //URL is the  hot page of the subreddit
	Hot           string `json:"hot"`
	New           string `json:"new"`
	Top           string `json:"top"`
	Rising        string `json:"rising"`
	Controversial string `json:"controversial"`
}

var err error

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"

//wg wait group for all downloads to complete
var wg sync.WaitGroup

//ScrapeImages downloads #num images from the given subreddit tab
//into a new folder at path
func (r SubReddit) ScrapeImages(tab string, num uint8, path string, linksFrom string) uint8 {

	err = os.MkdirAll(path, os.ModePerm)
	HandleErr(err)
	var baseURL string
	switch tab {
	case "top":
		baseURL = r.Top + "/?t=" + linksFrom
	}
	LinkCrawl(baseURL)

	return num + 1

}

var baseURL string

//LinkCrawl finds the required links
func LinkCrawl(url string) {
	baseURL = url
	linkCrawler(url)
}

func linkCrawler(url string) {
	req, err := http.NewRequest("GET", url, nil)
	HandleErr(err)
	req.Header.Set("User-Agent", userAgent)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	HandleErr(err)

	//finds the data urls from the entire page
	links := doc.Find("#siteTable div").Map(func(index int, item *goquery.Selection) string {
		link, _ := item.Attr("data-url")
		return link
	})

	//finds url of the next button
	var next string
	doc.Find("span .next-button a").Each(func(index int, item *goquery.Selection) {
		next, _ = item.Attr("href")
		fmt.Println(next)
	})

	//downloads each data url link if its a file
	for _, v := range links {
		if len(v) > 5 && v[len(v)-4:len(v)-3] == "." {
			wg.Add(1)
			go Download(v, "/Users/Reddy/gg")
		}
	}

	wg.Add(1)
	go linkCrawler(next)
	wg.Wait()
}

//HandleErr handles error
func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Download s from the url  and saves in path folder
func Download(url string, path string) {
	defer wg.Done()
	// req, err := http.NewRequest("GET", url, nil)
	// HandleErr(err)
	// req.Header.Set("User-Agent", userAgent)
	// c := &http.Client{}
	resp, err := http.Get(url)
	HandleErr(err)
	defer resp.Body.Close()
	tmp := strings.SplitAfter(url, "/")
	path = filepath.Join(path, tmp[len(tmp)-1])
	file, err := os.Create(path)
	defer file.Close()
	HandleErr(err)

	_, err = io.Copy(file, resp.Body)
	HandleErr(err)

}