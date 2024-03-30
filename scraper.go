package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gocolly/colly"
)

type HackathonEvent struct {
	thumbnailUrl, logoUrl, name, date, location, eventType string
}

func main() {

	c := colly.NewCollector()

	MLH_event_website := "https://mlh.io/seasons/2024/events"

	// an array to keep all hackathon event
	// var hackathonEvents []HackathonEvent
	hackathonEvents := []HackathonEvent{}

	// Testing purpose: print out the HTML text
	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println(string(r.Body))
	// })

	// print error when some is wrong
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("div.inner", func(e *colly.HTMLElement) {

		// <img src="" />
		// get the inner child attributes of src
		thumbnailUrl := e.ChildAttr("div.img-wrap img", "src")
		logoUrl := e.ChildAttr("div.event-logo img", "src")

		// get the text instead of the attributes
		name := e.ChildText("h3.event-name")
		date := e.ChildText("p.event-date")
		location := e.ChildText("div.event-location")
		eventType := e.ChildText("div.event-hybrid-notes")

		hackathonEvent := HackathonEvent{
			thumbnailUrl: thumbnailUrl,
			logoUrl:      logoUrl,
			name:         name,
			date:         date,
			location:     location,
			eventType:    eventType,
		}

		hackathonEvents = append(hackathonEvents, hackathonEvent)
	})

	// Create Proxy API URL, Simlar to URL in JS()
	websiteURL, err := url.Parse("https://proxy.scrapeops.io/v1/")
	if err != nil {
		log.Fatal(err)
	}

	// Add Query Parameters (api_key and url) and perform URL encoding such as adding %20
	q := websiteURL.Query()
	q.Set("api_key", "8a49e4a8-103a-41cc-9a13-e6afe2004745")
	q.Set("url", MLH_event_website)
	websiteURL.RawQuery = q.Encode()

	// Request Page
	c.Visit(websiteURL.String())

	// opening the CSV file
	file, err := os.Create("hackathons.csv")
	if err != nil {
		log.Fatalln("Failed to create hackathons.csv file", err)
	}

	defer file.Close() // close the IO

	// add a writter to the csv file
	writer := csv.NewWriter(file)

	// define the header
	headers := []string{
		"thumbnailUrl",
		"logoUrl",
		"name",
		"date",
		"location",
		"eventType",
	}

	// add header to the file
	writer.Write(headers)

	// add each hackathon event to csv file
	// iterate over each hackathon
	for _, event := range hackathonEvents {

		instance := []string{
			event.thumbnailUrl,
			event.logoUrl,
			event.name,
			event.date,
			event.location,
			event.eventType,
		}

		// add it to the csv file
		writer.Write(instance)
	}

	defer writer.Flush() // clean up the writer buffered data

	fmt.Println(hackathonEvents[0])
}
