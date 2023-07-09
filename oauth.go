package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	REDIRECT_URI = "http://localhost:8080/callback"
	AUTH_URL     = "https://id.twitch.tv/oauth2/authorize"
	TOKEN_URL    = "https://id.twitch.tv/oauth2/token"
)

var codeChan = make(chan string, 1)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func authorize() error {
	// Redirect the user to the Twitch authorization page
	fmt.Printf("Visit this link to authorize:\n%s\n", fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		AUTH_URL,
		CLIENT_ID,
		REDIRECT_URI,
		url.QueryEscape("chat:read chat:edit")))

	// Set up a temporary HTTP server to receive the authorization callback
	http.HandleFunc("/callback", handleCallback)
	go func() {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the authorization code and exchange it for an access token
	authCode := <-codeChan

	token, err := getToken(authCode)
	if err != nil {
		return fmt.Errorf("error exchanging auth code for token: %w", err)
	} else if token.AccessToken == "" {
		return fmt.Errorf("error getting token: empty")
	}

	ACCESS_TOKEN = token.AccessToken
	REFRESH_TOKEN = token.RefreshToken
	fmt.Println("token expires in", time.Duration(time.Duration(token.ExpiresIn)*time.Second).String())
	fmt.Println()

	close(codeChan)
	return nil
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code != "" {
		_, err := w.Write([]byte("Authorization code received. You can close this tab now."))
		if err != nil {
			log.Println("Failed to write response:", err)
		}
	} else {
		if _, err := w.Write([]byte("Failed to receive authorization code.")); err != nil {
			log.Println("Failed to write response:", err)
		}
	}

	go func() { codeChan <- code }()
}

func getToken(authCode string) (*tokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", CLIENT_ID)
	data.Set("client_secret", CLIENT_SECRET)
	data.Set("code", authCode)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", REDIRECT_URI)

	resp, err := http.PostForm(TOKEN_URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}
