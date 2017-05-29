package rscraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"

var wg sync.WaitGroup

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
					v = strings.Replace(v, "gallery", "download", -1)
					dlinks = append(dlinks, v)
				}
				if strings.Contains(v, "imgur.com/a/") {
					v = v + "/zip"
					dlinks = append(dlinks, v)
				}
				if strings.Contains(v, "imgur.com/download/") {
					v = strings.Replace(v, "/download", "", -1)
					dlinks = append(dlinks, v)
				}
			}
		}
	}
	return dlinks
}

//LinkCrawl finds the required links
func LinkCrawl(url string, path string) (err error) {
	wg.Add(1)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.Chdir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	linkCrawler(url)
	return err
}

func linkCrawler(url string) (err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//finds the data urls from the entire page
	links := doc.Find("#siteTable div").Map(func(index int, item *goquery.Selection) string {
		link, _ := item.Attr("data-url")
		return link
	})

	links = downloadbleLinks(links)
	println(len(links))
	for _, v := range links {
		wg.Add(1)
		go Download(v)
	}

	//finds url of the next button
	var next string
	doc.Find("span .next-button a").Each(func(index int, item *goquery.Selection) {
		next, _ = item.Attr("href")
		fmt.Println(next)
	})

	if len(next) > 0 {
		go linkCrawler(next)
	} else {
		wg.Done()
		println("crawling done")
	}
	wg.Wait()

	return err
}

//Download saves the contents from url into file
func Download(url string) (err error) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	contentType := resp.Header.Get("Content-Type")
	filename := ""
	if contentType == "application/zip" {
		filename = resp.Header.Get("Content-Disposition")
		//error may occur if split doesnt zero length slice
		filename = strings.Split(filename, "\"")[1]
	}
	if strings.Contains(contentType, "image") {
		filename = resp.Header.Get("Content-Disposition")
		if filename == "" {
			tmp := strings.Split(url, "/")
			filename = tmp[len(tmp)-1]
		} else {
			filename = strings.Split(filename, "\"")[1]
		}

	}
	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			log.Println(err)
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return err
}
