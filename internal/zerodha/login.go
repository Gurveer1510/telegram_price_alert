package zerodha

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/pquerna/otp/totp"
)

type ZerodhaClient struct {
	AccessToken string
	ApiKey      string
	ApiSecret   string
	Username    string
	Password    string
	TotpSecret  string
}

type loginResponse struct {
	Data struct {
		RequestID string `json:"request_id"`
	} `json:"data"`
}

func NewZerodhaClient(apiKey, username, password, totpSecret, apiSecret string) *ZerodhaClient {
	zclient := &ZerodhaClient{
		ApiKey:     apiKey,
		Username:   username,
		Password:   password,
		TotpSecret: totpSecret,
		ApiSecret: apiSecret,
	}
	accessToken, _ := zclient.GetAccessToken()
	zclient.AccessToken = accessToken
	return zclient
}

func (z *ZerodhaClient) GetAccessToken() (string, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Query().Get("request_token") != "" {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	requestID, err := z.login(client)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	err = z.submiTOTP(client, requestID)
	if err != nil {
		return "", fmt.Errorf("submiting totp failed: %w", err)
	}

	requestToken, err := z.getRequestToken(client)
	if err != nil {
		return "", fmt.Errorf("getting request token failed: %w", err)
	}

	// Step 4: Exchange for access_token
	accessToken, err := z.generateSession(requestToken)
	if err != nil {
		return "", fmt.Errorf("generating session failed: %w", err)
	}
	return accessToken, nil
}

func (z *ZerodhaClient) login(client *http.Client) (string, error) {

	urlString := "https://kite.zerodha.com/api/login"
	val := url.Values{
		"user_id":  {z.Username},
		"password": {z.Password},
	}

	resp, err := client.PostForm(urlString, val)
	if err != nil {
		log.Println(err)
	}

	var result loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.RequestID, nil
}

func (z *ZerodhaClient) submiTOTP(client *http.Client, requestID string) error {
	code, err := totp.GenerateCode(z.TotpSecret, time.Now())
	if err != nil {
		return err
	}

	resp, err := client.PostForm("https://kite.zerodha.com/api/twofa", url.Values{
		"user_id":     {z.Username},
		"request_id":  {requestID},
		"twofa_value": {code},
		"twofa_type":  {"totp"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}

func (z *ZerodhaClient) getRequestToken(client *http.Client) (string, error) {
	loginURL := fmt.Sprintf(
		"https://kite.trade/connect/login?api_key=%s&v=3", z.ApiKey,
	)

	resp, err := client.Get(loginURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	log.Printf("[debug] status=%d request_url=%s location=%s",
		resp.StatusCode, resp.Request.URL, location)

	var finalURL *url.URL

	if location != "" {
		// Use resp.Request.URL as base to resolve relative redirects
		finalURL, err = resp.Request.URL.Parse(location)
		if err != nil {
			return "", fmt.Errorf("parsing location header: %w", err)
		}
	} else {
		// No Location header — we must have followed all the way through
		finalURL = resp.Request.URL
	}

	params, err := url.ParseQuery(finalURL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("parsing query: %w", err)
	}

	token := params.Get("request_token")
	log.Printf("[debug] raw query: %s", finalURL.RawQuery)
	log.Printf("[debug] request_token: %s", token)

	if token == "" {
		return "", fmt.Errorf("request_token not found in: %s", finalURL)
	}
	return token, nil
}

func (z *ZerodhaClient) generateSession(requestToken string) (string, error) {
	// Checksum = sha256(api_key + request_token + api_secret)
	h := sha256.New()
	h.Write([]byte(z.ApiKey + requestToken + z.ApiSecret))
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	resp, err := http.PostForm("https://api.kite.trade/session/token", url.Values{
		"api_key":       {z.ApiKey},
		"request_token": {requestToken},
		"checksum":      {checksum},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Data.AccessToken == "" {
		return "", fmt.Errorf("empty access token in response")
	}
	return result.Data.AccessToken, nil
}
