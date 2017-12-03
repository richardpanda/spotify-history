package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

type Artist struct {
	Name string `json:"name"`
}

type PlayHistory struct {
	Track Track `json:"track"`
}

type ResponseBody struct {
	Items []PlayHistory `json:"items"`
}

type Track struct {
	Artists []Artist `json:"artists"`
	Name    string   `json:"name"`
}

func main() {
	token := os.Getenv("OAUTH_TOKEN")
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played", nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 401 {
		log.Fatal(string(body))
	}

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Println("\n")
	fmt.Fprintln(w, "Title\tArtist\t")
	for _, playHistory := range responseBody.Items {
		trackName := playHistory.Track.Name
		artists := playHistory.Track.Artists

		artistNames := make([]string, len(artists))
		for i, artist := range artists {
			artistNames[i] = artist.Name
		}

		fmt.Fprintf(w, "%s\t%s\n", trackName, strings.Join(artistNames, ", "))
	}
	w.Flush()
}
