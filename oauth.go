package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	ClientId     string
	ClientSecret string
)

const (
	RedirectURI = "http://localhost:8080/callback"
	AuthURL     = "https://id.twitch.tv/oauth2/authorize"
	TokenURL    = "https://id.twitch.tv/oauth2/token"
)

var codeChan = make(chan string)

func oauth() string {
	ClientId = os.Getenv("GOTATO_CLIENT_ID")
	ClientSecret = os.Getenv("GOTATO_CLIENT_SECRET")
	// Step 1: Redirect the user to the Twitch authorization page
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=channel:read:subscriptions", AuthURL, ClientId, RedirectURI)
	fmt.Printf("Please visit the following URL to authorize the application:\n%s\n\n", authURL)

	// Step 2: Set up a temporary HTTP server to receive the authorization callback
	http.HandleFunc("/callback", handleCallback)
	go func() {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Step 3: Wait for the authorization code and exchange it for an access token
	fmt.Print("Enter the authorization code: ")
	authorizationCode := <-codeChan

	token, err := exchangeAuthorizationCode(authorizationCode)
	if err != nil {
		log.Fatal("Failed to exchange authorization code for token:", err)
	}

	// check if we received a token
	if token.AccessToken == "" {
		log.Fatal("Failed to exchange authorization code for token: empty token")
	}

	fmt.Println("\nAccess Token:", token.AccessToken)
	fmt.Println("Expires In:", token.ExpiresIn)
	fmt.Println("Refresh Token:", token.RefreshToken)

	return token.AccessToken
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code != "" {
		_, err := w.Write([]byte("Authorization code received. You can close this tab now."))
		if err != nil {
			log.Println("Failed to write response:", err)
		}
	} else {
		_, err := w.Write([]byte("Failed to receive authorization code."))
		if err != nil {
			log.Println("Failed to write response:", err)
		}
	}
	go func() { codeChan <- code }()
}

func exchangeAuthorizationCode(authorizationCode string) (*OAuthToken, error) {
	data := url.Values{}
	data.Set("client_id", ClientId)
	data.Set("client_secret", ClientSecret)
	data.Set("code", authorizationCode)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", RedirectURI)

	resp, err := http.PostForm(TokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	token := &OAuthToken{}
	if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
