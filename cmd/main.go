package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/richardpanda/spotify-history/spotify"
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

func main() {
	clientID := "788fcd8f8c484f6e9df10ffab16c8844"
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := "http://localhost:7000/"

	err := openAuthWindow(clientID, redirectURI)
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

		fmt.Printf("%+v\n", auth)
	})

	log.Fatal(http.ListenAndServe(":7000", nil))

	// token := os.Getenv("OAUTH_TOKEN")
	// client := &http.Client{}

	// req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if resp.StatusCode == 401 {
	// 	log.Fatal(string(body))
	// }

	// var responseBody ResponseBody
	// err = json.Unmarshal(body, &responseBody)

	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	// fmt.Println("\n")
	// fmt.Fprintln(w, "Title\tArtist\t")
	// for _, playHistory := range responseBody.Items {
	// 	trackName := playHistory.Track.Name
	// 	artists := playHistory.Track.Artists

	// 	artistNames := make([]string, len(artists))
	// 	for i, artist := range artists {
	// 		artistNames[i] = artist.Name
	// 	}

	// 	fmt.Fprintf(w, "%s\t%s\n", trackName, strings.Join(artistNames, ", "))
	// }
	// w.Flush()
}
