package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

// type HackathonEvent struct {
// 	imageUrl, logoUrl, name, date, location, notes string
// }

func main() {

	// var hackathonEvents []HackathonEvent

	c := colly.NewCollector()

	MLH_event_website := "https://mlh.io/seasons/2024/events"

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 ")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	// print when there is error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		// hackathonEvent = HackathonEvent{}

		// imageUrl := e.ChildAttr("img", "src")

		// fmt.Println(imageUrl)
		fmt.Println(e.Text)
	})

	c.Visit(MLH_event_website)
}
