package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Making structures to read JSON response from LastFM API
type Image struct {
	Small string `json:"#text"`
}

type TrackMatches struct {
	Tracks []Track `json:"track"`
}

type SearchResults struct {
	TrackMatches TrackMatches `json:"results"`
}
type Track struct {
	Name   string  `json:"name"`
	Artist string  `json:"artist"`
	Url    string  `json:"url"`
	Image  []Image `json:"image"`
}

func SearchName(query string) []Track {
	// Asigning API_KEY var with opening .env file
	API_KEY, exists := os.LookupEnv("LASTFM_TOKEN")
	if !exists {
		log.Printf("BOT_TOKEN not found in environment variables")
	}

	// Making request to API with formating our query
	query = strings.ReplaceAll(query, " ", "+")
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.search&track=%s&api_key=%s&format=json&limit=5", query, API_KEY)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// Closing our request with end of function
	defer response.Body.Close()

	// Converting JSON response to code structure
	var data []byte
	data, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var searchResults struct {
		Results struct {
			TrackMatches TrackMatches `json:"trackmatches"`
		} `json:"results"`
	}

	err = json.Unmarshal(data, &searchResults)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}
	if len(searchResults.Results.TrackMatches.Tracks) == 0 {
		fmt.Println("No tracks found.")
	}

	return searchResults.Results.TrackMatches.Tracks
}

// Also making structures to convert JSON response
type Artist struct {
	Name string `json:"name"`
}

type TopTrack struct {
	Name   string `json:"name"`
	Artist Artist `json:"artist"`
	URL    string `json:"url"`
}

type TopList struct {
	Tracks []TopTrack `json:"track"`
}

func GetTopTracks() []TopTrack {
	// Asinging variable using .env file
	API_KEY, exists := os.LookupEnv("LASTFM_TOKEN")
	if !exists {
		log.Printf("BOT_TOKEN not found in environment variables")
	}

	// Making HTTP request to LastFM API
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=chart.gettoptracks&api_key=%s&format=json", API_KEY)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error making GET request:", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", response.StatusCode)
	}

	// Reading JSON response and converting it to needed structure
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var searchResults struct {
		TopTracks struct {
			Tracks []TopTrack `json:"track"`
		} `json:"tracks"`
	}

	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	if len(searchResults.TopTracks.Tracks) == 0 {
		fmt.Println("No tracks found.")
	}

	return searchResults.TopTracks.Tracks
}
