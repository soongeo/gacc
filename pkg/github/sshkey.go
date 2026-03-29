package github

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// AddSSHPublicKey GitHub API를 통해 생성된 공개키를 업로드합니다.
func AddSSHPublicKey(accessToken, accountName, publicKey string) error {
	client := resty.New()

	body := map[string]string{
		"title": fmt.Sprintf("gacc-%s", accountName),
		"key":   publicKey,
	}

	resp, err := client.R().
		SetAuthToken(accessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetBody(body).
		Post("https://api.github.com/user/keys")

	if err != nil {
		return fmt.Errorf("network error during SSH key upload: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to register SSH key to GitHub (Status: %v): %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// SSHKey GitHub API에서 반환하는 SSH 키 구조체
type SSHKey struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// DeleteSSHPublicKey GitHub 계정에 등록된 ssh 키 중 title이 일치하는 키를 찾아 삭제합니다.
func DeleteSSHPublicKey(accessToken, accountName string) error {
	client := resty.New()
	targetTitle := fmt.Sprintf("gacc-%s", accountName)

	// 1. 등록된 SSH 키 목록 가져오기
	resp, err := client.R().
		SetAuthToken(accessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		Get("https://api.github.com/user/keys")

	if err != nil {
		return fmt.Errorf("network error fetching SSH key list: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to fetch SSH key list (Status: %v)", resp.StatusCode())
	}

	var keys []SSHKey
	if err := json.Unmarshal(resp.Body(), &keys); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// 2. 일치하는 키 찾기
	var targetKeyID int
	for _, k := range keys {
		if k.Title == targetTitle {
			targetKeyID = k.ID
			break
		}
	}

	// 일치하는 키가 없으면 그냥 반환
	if targetKeyID == 0 {
		return nil
	}

	// 3. 키 삭제 요청
	delResp, err := client.R().
		SetAuthToken(accessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		Delete(fmt.Sprintf("https://api.github.com/user/keys/%d", targetKeyID))

	if err != nil {
		return fmt.Errorf("network error deleting SSH key: %w", err)
	}

	if delResp.IsError() {
		return fmt.Errorf("failed to delete GitHub SSH key (Status: %v)", delResp.StatusCode())
	}

	return nil
}
