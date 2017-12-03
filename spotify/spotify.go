package spotify

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

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
