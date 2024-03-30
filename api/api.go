package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// must be capitalized to be exported for the keys
type HackathonEvent struct {
	ThumbnailUrl string `json:"thumbnailUrl"`
	LogoUrl      string `json:"logoUrl"`
	Name         string `json:"name"`
	Date         string `json:"date"`
	Location     string `json:"location"`
	EventType    string `json:"eventType"`
}

func readFromCSV(fileName string) []HackathonEvent {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	// Closes the file
	defer file.Close()

	// create a csv reader file object to the read the file
	reader := csv.NewReader(file)

	// read all the records from the csv file
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	hackathonsEvents := []HackathonEvent{}

	// start form the second csv header
	for _, record := range records[1:] {
		hackathonsEvents = append(hackathonsEvents, HackathonEvent{
			ThumbnailUrl: record[0],
			LogoUrl:      record[1],
			Name:         record[2],
			Date:         record[3],
			Location:     record[4],
			EventType:    record[5],
		})
	}

	return hackathonsEvents
}

func getAllHackathonEvents(c *gin.Context) {

	hackathonEvents := readFromCSV("../hackathons_2024.csv")

	c.IndentedJSON(200, hackathonEvents)

}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	router.GET("/api/hackathons", getAllHackathonEvents)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalln("Failed to start the server", err)
	}

}
