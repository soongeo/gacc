package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type GitHubUser struct {
	Name  string `json:"name"`
	Login string `json:"login"`
}

type GitHubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

// UserInfo represents the fetched account identity
type UserInfo struct {
	Name  string
	Email string
}

// GetUserInfo fetches the authenticated user's name and primary email.
func GetUserInfo(accessToken string) (UserInfo, error) {
	client := resty.New()

	// 1. Fetch User Profile
	resp, err := client.R().
		SetAuthToken(accessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		Get("https://api.github.com/user")

	if err != nil {
		return UserInfo{}, fmt.Errorf("network error fetching user profile: %w", err)
	}

	if resp.IsError() {
		return UserInfo{}, fmt.Errorf("failed to fetch user profile (Status: %v)", resp.StatusCode())
	}

	var user GitHubUser
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return UserInfo{}, fmt.Errorf("failed to parse user response: %w", err)
	}

	displayName := user.Name
	if strings.TrimSpace(displayName) == "" {
		displayName = user.Login
	}

	// 2. Fetch User Emails
	emailResp, err := client.R().
		SetAuthToken(accessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		Get("https://api.github.com/user/emails")

	if err != nil {
		return UserInfo{}, fmt.Errorf("network error fetching user emails: %w", err)
	}

	var email string
	if !emailResp.IsError() {
		var emails []GitHubEmail
		if err := json.Unmarshal(emailResp.Body(), &emails); err == nil {
			for _, e := range emails {
				if e.Primary {
					email = e.Email
					break
				}
			}
			// Fallback: If no primary marked, pick the first one
			if email == "" && len(emails) > 0 {
				email = emails[0].Email
			}
		}
	}

	return UserInfo{
		Name:  displayName,
		Email: email,
	}, nil
}
