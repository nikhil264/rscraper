package rscraper

import "testing"

func TestDownload(t *testing.T) {
	url := "http://imgur.com/AdpTzNI.png"
	err := Download(url)
	if err != nil {
		t.Error("failed direct image case", err)
	}

	url = "http://imgur.com/download/AdpTzNI"
	err = Download(url)
	if err != nil {
		t.Error("failed single image case", err)
	}
	url = "http://imgur.com/a/jfLYQ/zip"
	err = Download(url)
	if err != nil {
		t.Error("failed album image case", err)
	}
}

func TestDownloadbleLinks(t *testing.T) {

	url := "https://i.redd.it/yan6sp65qwcy.jpg"
	if downloadbleLinks([]string{url})[0] != url {
		t.Error("failed direct image case")
	}

	url = "http://imgur.com/gallery/FUWRGYP"
	if downloadbleLinks([]string{url})[0] != "http://imgur.com/a/FUWRGYP/zip" {
		t.Error("failed zip  case")
	}
	if downloadbleLinks([]string{url})[1] != "http://imgur.com/download/FUWRGYP" {
		t.Error("failed single image case")
	}
}
