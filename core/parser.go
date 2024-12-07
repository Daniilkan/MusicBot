package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

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

const API_KEY string = ""

func SearchName(query string) []Track {

	query = strings.ReplaceAll(query, " ", "+")
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.search&track=%s&api_key=%s&format=json&limit=5", query, API_KEY)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

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
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=chart.gettoptracks&api_key=%s&format=json", API_KEY)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error making GET request:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", response.StatusCode)
	}

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
