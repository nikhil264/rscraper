package rscraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

//LinkCrawl finds the required links
func LinkCrawl(baseURL string) {
	doc, err := goquery.NewDocument(baseURL)
	HandleErr(err)

	doc.Find(".siteTable").Each(func(index int, item *goquery.Selection) {
		linkTag := item.Find("div")
		after, _ := linkTag.Attr("data-fullname")
		link, _ := linkTag.Attr("data-url")
		fmt.Printf("%s  %s\n", after, link)
	})
}

//HandleErr handles error
func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Download s from the url  and saves in path folder
func Download(url string, path string) {
	r, err := http.Get(url)
	HandleErr(err)
	defer r.Body.Close()
	tmp := strings.SplitAfter(url, "/")
	path = filepath.Join(path, tmp[len(tmp)-1])
	file, err := os.Create(path)
	defer file.Close()
	HandleErr(err)

	_, err = io.Copy(file, r.Body)
	HandleErr(err)

}
