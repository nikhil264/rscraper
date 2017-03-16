package main

import "github.com/nikhil264/rscraper"

func main() {
	baseURL := "https://www.reddit.com/r/wallpaper/top/?sort=top&t=week"
	rscraper.LinkCrawl(baseURL)
}
