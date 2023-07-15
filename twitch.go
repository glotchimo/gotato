package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/pkg/browser"
)

const (
	REDIRECT_URI = "http://localhost:8080/callback"
	AUTH_URL     = "https://id.twitch.tv/oauth2/authorize"
	TOKEN_URL    = "https://id.twitch.tv/oauth2/token"
	SCOPES       = "chat:read chat:edit moderator:manage:banned_users user:manage:whispers"
)

var codeChan = make(chan string, 1)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func authenticate() error {
	// Redirect the user to the Twitch authorization page
	browser.OpenURL(fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		AUTH_URL,
		CLIENT_ID,
		REDIRECT_URI,
		url.QueryEscape(SCOPES)))

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

	close(codeChan)
	return nil
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code != "" {
		_, err := w.Write([]byte("Got it! Head back to your terminal."))
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

func refreshToken() error {
	params := url.Values{}
	params.Add("grant_type", `refresh_token`)
	params.Add("refresh_token", REFRESH_TOKEN)
	params.Add("client_id", CLIENT_ID)
	params.Add("client_secret", CLIENT_SECRET)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", body)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	var token tokenResponse
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return err
	}

	ACCESS_TOKEN = token.AccessToken
	REFRESH_TOKEN = token.RefreshToken

	return nil
}

func createAPIClient() error {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        CLIENT_ID,
		UserAccessToken: ACCESS_TOKEN,
	})
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	CLIENT_API = client
	return nil
}

func timeout(id string) error {
	if _, err := CLIENT_API.BanUser(&helix.BanUserParams{
		BroadcasterID: BROADCASTER_ID,
		ModeratorId:   BROADCASTER_ID,
		Body: helix.BanUserRequestBody{
			UserId:   id,
			Duration: TIMEOUT_DURATION,
			Reason:   "lost to potato",
		},
	}); err != nil {
		return fmt.Errorf("error timing out loser: %w", err)
	}

	return nil
}

func whisper(id string, message string) error {
	if _, err := CLIENT_API.SendUserWhisper(&helix.SendUserWhisperParams{
		FromUserID: BROADCASTER_ID,
		ToUserID:   id,
		Message:    message,
	}); err != nil {
		return fmt.Errorf("error sending whisper: %w", err)
	}

	return nil
}
