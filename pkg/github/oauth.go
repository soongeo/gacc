package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

// DeviceCodeResponse GitHub Device Code 요청의 응답 구조
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// AccessTokenResponse access token 요쳥의 응답 구조
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

// openBrowser OS 기본 브라우저를 열어줍니다.
func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

// StartDeviceFlow GitHub Device Flow를 시작하고 Access Token을 반환합니다.
func StartDeviceFlow(clientID string) (string, error) {
	if clientID == "" {
		return "", errors.New("GitHub Client ID is required. Please check your config file or GACC_GITHUB_CLIENT_ID environment variable.")
	}

	client := resty.New()

	// 1. Device Code 요청
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"client_id": clientID,
			"scope":     "admin:public_key user:email",
		}).
		Post("https://github.com/login/device/code")

	if err != nil {
		return "", fmt.Errorf("failed to request device code: %w", err)
	}

	var deviceResp DeviceCodeResponse
	if err := json.Unmarshal(resp.Body(), &deviceResp); err != nil {
		return "", fmt.Errorf("failed to parse device code response: %w", err)
	}

	// 2. Guide user with code and URL
	fmt.Println("\n=======================================================")
	fmt.Printf("🔑 Please enter the following code in your browser: %s\n", deviceResp.UserCode)
	fmt.Println("=======================================================")
	
	// 브라우저 자동 실행
	if err := openBrowser(deviceResp.VerificationURI); err != nil {
		fmt.Printf("⚠️ Cannot open browser automatically. Please visit the following URL manually:\n   👉 %s\n\n", deviceResp.VerificationURI)
	} else {
		fmt.Println("🌐 Browser opened automatically.")
	}

	fmt.Println("⏳ Waiting for the authentication process to complete...")

	// 3. Access Token 폴링
	interval := time.Duration(deviceResp.Interval) * time.Second
	for {
		time.Sleep(interval)

		tokenResp, err := client.R().
			SetHeader("Accept", "application/json").
			SetFormData(map[string]string{
				"client_id":   clientID,
				"device_code": deviceResp.DeviceCode,
				"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
			}).
			Post("https://github.com/login/oauth/access_token")

		if err != nil {
			return "", fmt.Errorf("access token request error: %w", err)
		}

		var tokenData AccessTokenResponse
		if err := json.Unmarshal(tokenResp.Body(), &tokenData); err != nil {
			return "", fmt.Errorf("failed to parse access token response: %w", err)
		}

		// 성공적으로 토큰을 받음
		if tokenData.AccessToken != "" {
			return tokenData.AccessToken, nil
		}

		// 인증 진행 중
		if tokenData.Error == "authorization_pending" {
			continue
		} else if tokenData.Error == "slow_down" {
			interval += 5 * time.Second // 추가 대기
			continue
		} else if tokenData.Error != "" {
			return "", fmt.Errorf("GitHub Auth Error: %s", tokenData.Error)
		}
	}
}
