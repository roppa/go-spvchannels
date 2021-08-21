package spvchannels

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ClientConfig struct {
	Insecure bool // equivalent curl -k
	BaseURL  string
	Version  string
	User     string
	Passwd   string
	Token    string
}

type Client struct {
	cfg        ClientConfig
	HTTPClient *http.Client
}

func NewClient(c ClientConfig) *Client {
	var httpClient http.Client
	if c.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = http.Client{
			Timeout:   time.Minute,
			Transport: tr,
		}
	} else {
		httpClient = http.Client{
			Timeout: time.Minute,
		}
	}

	return &Client{
		cfg:        c,
		HTTPClient: &httpClient,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type GetChannelRequest struct {
	Method    string `json:"method"`
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type GetChannelReply struct {
	Id          string `json:"id"`
	Href        string `json:"href"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Locked      bool   `json:"locked"`
	Head        int    `json:"head"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
	AccessTokens []struct {
		Id          string `json:"id"`
		Token       string `json:"token"`
		Description string `json:"description"`
		CanRead     bool   `json:"can_read"`
		CanWrite    bool   `json:"can_write"`
	} `json:"access_tokens"`
}

func (c *Client) GetChannel(ctx context.Context, r GetChannelRequest) (*GetChannelReply, error) {
	req, err := http.NewRequestWithContext(ctx, r.Method, fmt.Sprintf("https://%s/api/%s/account/%s/channel/%s", c.cfg.BaseURL, c.cfg.Version, r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	res := GetChannelReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	if c.cfg.Token == "" {
		req.SetBasicAuth(c.cfg.User, c.cfg.Passwd)
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.Token))
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := successResponse{
		Code: res.StatusCode,
		Data: v,
	}

	if err = json.NewDecoder(res.Body).Decode(&fullResponse.Data); err != nil {
		return err
	}

	return nil
}
