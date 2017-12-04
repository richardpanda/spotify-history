package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Artist struct {
	Name string `json:"name"`
}

type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRequestParams struct {
	ClientID, ClientSecret string
	Code                   string
	RedirectURI            string
	RefreshToken           string
}

type PlayHistory struct {
	Track Track `json:"track"`
}

type Track struct {
	Artists []Artist `json:"artists"`
	Name    string   `json:"name"`
}

type TracksResponseBody struct {
	Items []PlayHistory `json:"items"`
}

func FetchAuth(params AuthRequestParams) (*Auth, error) {
	var body url.Values
	if params.RefreshToken == "" {
		body = url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {params.Code},
			"redirect_uri": {params.RedirectURI},
		}
	} else {
		body = url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {params.RefreshToken},
		}
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	encodedString := base64.StdEncoding.EncodeToString([]byte(params.ClientID + ":" + params.ClientSecret))
	req.Header.Set("Authorization", "Basic "+encodedString)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var auth Auth
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func RecentlyPlayedTracks(accessToken string) ([]Track, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}

	respBody := &TracksResponseBody{}
	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return nil, err
	}

	tracks := make([]Track, len(respBody.Items))
	for i, item := range respBody.Items {
		tracks[i] = item.Track
	}
	return tracks, nil
}
