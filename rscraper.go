package rscraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"

var downloadsWg sync.WaitGroup

var n int // no of remaining downloads.
var total int
var baseURL string

func downloadbleLinks(links []string) (dlinks []string) {
	//downloads each data url link if its a file
	for _, v := range links {
		if len(v) > 5 {
			if v[len(v)-4:len(v)-3] == "." {
				dlinks = append(dlinks, v)
			} else {
				if strings.Contains(v, "imgur.com/gallery/") {
					t := strings.Replace(v, "gallery", "a", -1) + "/zip"
					dlinks = append(dlinks, t)
					t = strings.Replace(v, "gallery", "download", -1)
					dlinks = append(dlinks, t)
				}
				if strings.Contains(v, "imgur.com/a/") {
					t := v + "/zip"
					dlinks = append(dlinks, t)
				}
				if strings.Contains(v, "imgur.com/download/") {
					t := strings.Replace(v, "/download", "", -1)
					dlinks = append(dlinks, t)
				}
			}
		}
	}
	return dlinks
}

//LinkCrawl finds the required links and creates a folder to store the downloaded files
func LinkCrawl(url string, path string) {
	baseURL = url
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)

	//intially no downloads are remaining
	n = 0
	total = 0
	os.MkdirAll(path, os.ModePerm)
	os.Chdir(path)

	linkCrawler(url)
}

func linkCrawler(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", userAgent)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(resp)

	//finds the data urls from the entire page
	links := doc.Find("#siteTable div").Map(func(index int, item *goquery.Selection) string {
		link, _ := item.Attr("data-url")
		return link
	})

	links = downloadbleLinks(links)

	for _, v := range links {
		downloadsWg.Add(1)
		n++
		// go Download(v)
		go fake(v)
	}

	//finds url of the next button
	var next string
	doc.Find("span .next-button a").Each(func(index int, item *goquery.Selection) {
		next, _ = item.Attr("href")
		fmt.Println(strings.TrimPrefix(next, baseURL))
		log.Printf(next)
	})

	for n >= 0 && len(next) > 0 {
		if n < 101 {
			linkCrawler(next)
		} else {
			time.Sleep(time.Second * 5)
			fmt.Println("Downloaded ", total, " items.")
		}
	}
	downloadsWg.Wait()
	println(n)
	println("crawling done Downloaded ", total, " items.")
}

//Download saves the contents from url into file
func Download(url string) (err error) {
	defer downloadsWg.Done()

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	contentType := resp.Header.Get("Content-Type")

	filename := ""
	if strings.Contains(contentType, "zip") || strings.Contains(contentType, "image") {
		filename = resp.Header.Get("Content-Disposition")
		//error may occur if split doesnt zero length slice
		if filename == "" {
			tmp := strings.Split(url, "/")
			filename = tmp[len(tmp)-1]
		} else {
			filename = strings.Split(filename, "\"")[1]
		}

	}

	if filename != "" {
		file, _ := os.Create(filename + ".part")
		defer file.Close()

		io.Copy(file, resp.Body)
		os.Rename(filename+".part", filename)
		log.Printf(filename + "  " + url)
	}
	n--
	total++
	return err
}

func fake(url string) {
	defer downloadsWg.Done()
	n--
	total++
}
