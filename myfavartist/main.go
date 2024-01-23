package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"net/url"
	// "os"
	//"strings"

	"github.com/gin-gonic/gin"
)

/*
{
  "lastfm_user": "",
  "lastfm_api_key": "",
  "musixmatch_api_key": "",
  "default_region": "us",
  "lastfm_api_url": "",
  "musixmatch_api_url": ""
}*/

// Configuration
type Config struct {
	LastFMUser       string `json:"lastfm_user"`
	LastFMAPIKey     string `json:"lastfm_api_key"`
	MusixMatchAPIKey string `json:"musixmatch_api_key"`
	DefaultRegion    string `json:"default_region"`
	LastFMAPIURL     string `json:"lastfm_api_url"`
	MusixMatchAPIURL string `json:"musixmatch_api_url"`
}

var config Config

func loadConfig() {
	jsonData, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var localConfig Config

	err = json.Unmarshal(jsonData, &localConfig)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	config = localConfig
	// Print the populated struct
	fmt.Printf("LastFM User: %s\nLastFM API Key: %s\nMusixMatch API Key: %s\nDefault Region: %s\nLastFM API URL: %s\nMusixMatch API URL: %s\n",
		config.LastFMUser, config.LastFMAPIKey, config.MusixMatchAPIKey, config.DefaultRegion, config.LastFMAPIURL, config.MusixMatchAPIURL)
}

// LastFMTopTracksResponse represents the structure of the Last.fm top tracks response
type LastFMTopTracksResponse struct {
	Tracks struct {
		Track []LastFMTrack `json:"track"`
	} `json:"tracks"`
	Attr struct {
		Country    string `json:"country"`
		Page       string `json:"page"`
		PerPage    string `json:"perPage"`
		TotalPages string `json:"totalPages"`
		Total      string `json:"total"`
	} `json:"@attr"`
}

// LastFMTrack represents a track in the Last.fm top tracks response
type LastFMTrack struct {
	Name       string        `json:"name"`
	Duration   string        `json:"duration"`
	Listeners  string        `json:"listeners"`
	Mbid       string        `json:"mbid"`
	URL        string        `json:"url"`
	Streamable Streamable    `json:"streamable"`
	Artist     LastFMArtist  `json:"artist"`
	Images     []LastFMImage `json:"image"`
	Attr       struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
}

// Streamable represents the streamable information in the Last.fm top tracks response
type Streamable struct {
	Text      string `json:"#text"`
	Fulltrack string `json:"fulltrack"`
}

// LastFMArtist represents an artist in the Last.fm top tracks response
type LastFMArtist struct {
	Name string `json:"name"`
	Mbid string `json:"mbid"`
	URL  string `json:"url"`
}

// LastFMImage represents an image in the Last.fm top tracks response
type LastFMImage struct {
	Text string `json:"#text"`
	Size string `json:"size"`
}

func main() {
	loadConfig()
	router := gin.Default()

	api := router.Group("/api/v1")
	api.GET("/artist/:region", getArtistInfo)

	router.Run(":8080")
}

func getArtistInfo(c *gin.Context) {
	region := c.Param("region")
	fmt.Println("region from req:", region)
	// Get conig
	// lastFMAPIKey := config.lastFMAPIKey
	// musixMatchAPIKey := config.musixMatchAPIKey
	defaultRegion := config.DefaultRegion
	// lastFMAPIURL := config.lastFMAPIURL
	// musixMatchAPIURL := config.musixMatchAPIURL
	fmt.Println("default region from config:", defaultRegion)

	if region == "" {
		region = defaultRegion
	}

	fmt.Println("config data: ", config.LastFMAPIURL) //lastFMAPIKey, musixMatchAPIKey, defaultRegion, lastFMAPIURL, musixMatchAPIURL)
	// Get toptrack in the region from Last.fm
	lastFMTrack, err := getLastFMTopTrack(region)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(lastFMTrack)

	// Get lyrics of the track from Musixmatch
	/*
		lyrics, err := getMusixmatchLyrics(lastFMTrack.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	// only commercial api provides search based on track name or artist name

	c.JSON(http.StatusOK, gin.H{
		"top_track": lastFMTrack.Name,
		"lyrics":    "",
		"artist_info": gin.H{
			"name":      lastFMTrack.Artist.Name,
			"listeners": lastFMTrack.Listeners,
		},
		"artist_image": lastFMTrack.Images[0].Text,
	})
}

func getLastFMTopTrack(region string) (LastFMTrack, error) {
	response, err := http.Get(fmt.Sprintf(config.LastFMAPIURL, region, config.LastFMAPIKey))
	if err != nil {
		fmt.Println("Error: Unable to get response from lastfm api")
		return LastFMTrack{}, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: Unable to read response body")
		return LastFMTrack{}, err
	}

	var result LastFMTopTracksResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error: Unable to unmarshall response body")
		return LastFMTrack{}, err
	}

	// Extract and return relevant information
	if len(result.Tracks.Track) > 0 {
		fmt.Printf("record exists\n")
		return result.Tracks.Track[0], nil
	}

	return LastFMTrack{}, nil
}

func getMusixmatchLyrics(track string) (string, error) {
	queryParams := url.Values{}
	queryParams.Add("q_track", track)
	queryParams.Add("apikey", config.MusixMatchAPIKey)
	requestURL := config.MusixMatchAPIURL + "track.search?" + queryParams.Encode()

	// request url ?
	fmt.Println("Request url for musixmatch:", requestURL)
	response, err := http.Get(requestURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	trackList := result["message"].(map[string]interface{})["body"].(map[string]interface{})["track_list"].([]interface{})
	if len(trackList) > 0 {
		lyrics := trackList[0].(map[string]interface{})["track"].(map[string]interface{})["lyrics"].(map[string]interface{})["lyrics_body"].(string)
		return lyrics, nil
	}

	return "", fmt.Errorf("Lyrics not found for the track: %s", track)
}
