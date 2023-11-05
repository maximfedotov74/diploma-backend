package yandexauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type YandexAuth struct {
	clientId     string
	cleintSecret string
	httpCleint   http.Client
}

func NewYandexAuth(cleintId string, clientSecret string) *YandexAuth {
	return &YandexAuth{
		clientId:     cleintId,
		cleintSecret: clientSecret,
		httpCleint:   http.Client{},
	}
}

func (ya *YandexAuth) GetYandexLoginUrl() string {
	url := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=code&client_id=%s", ya.clientId)
	return url
}

type yandexAuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type YandexLoginResponse struct {
	Birthday      string `json:"birthday"`
	Avatar        string `json:"default_avatar_id"`
	Email         string `json:"default_email"`
	FirstName     string `json:"first_name"`
	IsAvatarEmpty bool   `json:"is_avatar_empty"`
	LastName      string `json:"last_name"`
	Sex           string `json:"sex"`
	RealName      string `json:"real_name"`
}

func (ya *YandexAuth) GetUserByCode(code string) (*YandexLoginResponse, error) {

	formData := url.Values{}

	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)

	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "https://oauth.yandex.ru/token", body)

	if err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", ya.clientId, ya.cleintSecret)))
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))

	response, err := ya.httpCleint.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	yandexResponse := yandexAuthResponse{}

	err = json.Unmarshal(data, &yandexResponse)
	if err != nil {
		return nil, err
	}

	userReq, err := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return nil, err
	}
	userReq.Header.Set("Authorization", fmt.Sprintf("OAuth %s", yandexResponse.AccessToken))

	userResp, err := ya.httpCleint.Do(userReq)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	userData, err := io.ReadAll(userResp.Body)

	if err != nil {
		return nil, err
	}

	yandexLoginResponse := YandexLoginResponse{}

	err = json.Unmarshal(userData, &yandexLoginResponse)
	if err != nil {
		return nil, err
	}

	return &yandexLoginResponse, nil
}
