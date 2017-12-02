package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Song struct {
	Artists []Artist `json:"artists"`
	Name    string   `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}

type PlayHistory struct {
	Song Song `json:"item"`
}

func main() {
	token := os.Getenv("OAUTH_TOKEN")
	client := &http.Client{}
	var currentSong, previousSong Song

	for {
		func() {
			defer time.Sleep(1 * time.Minute)

			req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == 400 {
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			if resp.StatusCode == 401 {
				log.Fatal(string(body))
			}

			var playHistory PlayHistory
			err = json.Unmarshal(body, &playHistory)
			if err != nil {
				log.Fatal(err)
			}

			currentSong = playHistory.Song
			if previousSong.Name == currentSong.Name {
				return
			}

			previousSong = currentSong
			artistNames := make([]string, len(currentSong.Artists))
			for i, artist := range currentSong.Artists {
				artistNames[i] = artist.Name
			}

			fmt.Printf("%s - %s\n", currentSong.Name, strings.Join(artistNames, ", "))
		}()
	}
}
