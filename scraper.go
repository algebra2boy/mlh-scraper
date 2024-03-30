package main

import (
	"fmt"
	"os"
	"time"

	"encoding/csv"
	"log"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type HackathonEvent struct {
	thumbnailUrl, logoUrl, name, date, city, state, eventType string
}

func scrape(APIKey string, hackathonEvents []HackathonEvent) []HackathonEvent {
	c := colly.NewCollector()

	MLH_event_website := "https://mlh.io/seasons/2024/events"

	// print error when some is wrong
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("div.inner", func(e *colly.HTMLElement) {

		// <img src="" />
		// get the inner child attributes of src
		thumbnailUrl := e.ChildAttr("div.image-wrap img", "src") // /events/splashes
		logoUrl := e.ChildAttr("div.event-logo img", "src")      // /events/logos

		// get the text instead of the attributes
		name := e.ChildText("h3.event-name")
		date := e.ChildText("p.event-date")

		city := e.ChildText("div.event-location > span[itemprop='city']")
		state := e.ChildText("div.event-location > span[itemprop='state']")

		eventType := e.ChildText("div.event-hybrid-notes")

		hackathonEvent := HackathonEvent{
			thumbnailUrl: thumbnailUrl,
			logoUrl:      logoUrl,
			name:         name,
			date:         date,
			city:         city,
			state:        state,
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
	q.Set("api_key", APIKey)
	q.Set("url", MLH_event_website)
	websiteURL.RawQuery = q.Encode()

	// Request Page
	c.Visit(websiteURL.String())

	return hackathonEvents
}

func saveToCSV(filename string, hackathonEvents []HackathonEvent) {
	// opening the CSV file
	file, err := os.Create(filename)
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
		"city",
		"state",
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
			event.city,
			event.state,
			event.eventType,
		}

		// add it to the csv file
		writer.Write(instance)
	}

	defer writer.Flush() // clean up the writer buffered data
}

func main() {

	fmt.Println("Start scraping the website")

	// Load the .env file
	godotenv.Load()

	APIKey := os.Getenv("Proxy_API_KEY")

	if APIKey == "" {
		log.Fatal("Proxy_API_KEY is not found")
	}

	year := time.Now().Year()

	fileName := fmt.Sprintf("hackathons_%d.csv", year)

	// an array to keep all hackathon event
	// var hackathonEvents []HackathonEvent
	hackathonEvents := []HackathonEvent{}

	hackathonEvents = scrape(APIKey, hackathonEvents)

	saveToCSV(fileName, hackathonEvents)

	fmt.Println("Finish scraping the website")
}
