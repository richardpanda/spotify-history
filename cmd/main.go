package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/richardpanda/spotify-history/spotify"
)

func openAuthWindow(clientID, redirectURI string) error {
	responseType := "code"
	scope := "user-read-recently-played"
	v := url.Values{}
	v.Set("client_id", clientID)
	v.Set("response_type", responseType)
	v.Set("redirect_uri", redirectURI)
	v.Set("scope", scope)
	authorizeURL := fmt.Sprintf("https://accounts.spotify.com/authorize/?%s", v.Encode())

	err := exec.Command("open", authorizeURL).Start()
	if err != nil {
		return err
	}
	return nil
}

func printTracks(tracks []spotify.Track) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Println("\n")
	fmt.Fprintln(w, "Title\tArtist\t")
	for _, track := range tracks {
		artistNames := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artistNames[i] = artist.Name
		}

		fmt.Fprintf(w, "%s\t%s\n", track.Name, strings.Join(artistNames, ", "))
	}
	w.Flush()
}

func saveAuth(auth *spotify.Auth) error {
	b, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("auth.json", b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	clientID := "788fcd8f8c484f6e9df10ffab16c8844"
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := "http://localhost:7000/"

	_, err := os.Stat("auth.json")
	if err != nil {
		err = openAuthWindow(clientID, redirectURI)
		if err != nil {
			log.Fatal(err)
		}

		http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			params := spotify.AuthRequestParams{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Code:         r.Form.Get("code"),
				RedirectURI:  redirectURI,
			}

			auth, err := spotify.FetchAuth(params)
			if err != nil {
				log.Fatal(err)
			}

			err = saveAuth(auth)
			if err != nil {
				log.Fatal(err)
			}

			tracks, err := spotify.RecentlyPlayedTracks(auth.AccessToken)
			if err != nil {
				log.Fatal(err)
			}

			printTracks(tracks)
			os.Exit(0)
		})

		log.Fatal(http.ListenAndServe(":7000", nil))
	}

	fileBytes, err := ioutil.ReadFile("auth.json")
	if err != nil {
		log.Fatal(err)
	}

	auth := &spotify.Auth{}
	err = json.Unmarshal(fileBytes, auth)
	if err != nil {
		log.Fatal(err)
	}

	tracks, err := spotify.RecentlyPlayedTracks(auth.AccessToken)
	if err != nil {
		params := spotify.AuthRequestParams{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RefreshToken: auth.RefreshToken,
		}
		newAuth, err := spotify.FetchAuth(params)
		if err != nil {
			log.Fatal(err)
		}

		auth.AccessToken = newAuth.AccessToken
		err = saveAuth(auth)
		if err != nil {
			log.Fatal(err)
		}

		tracks, err = spotify.RecentlyPlayedTracks(auth.AccessToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	printTracks(tracks)
}
