package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/third-middleware/pkg/const"
	"github.com/go-resty/resty/v2"
)

func GetAccessToken(clientID, clientSecret, code string) (map[string]interface{}, error) {
	var err error
	reqURL := fmt.Sprintf(
		"%s?client_id=%s&client_secret=%s&code=%s",
		constant.GithubGetAccessTokenURL, clientID, clientSecret, code,
	)
	socksProxy := os.Getenv("ENV_RECAPTCHA_REQUEST_PROXY")

	cli := resty.New()

	if socksProxy != "" {
		cli = cli.SetProxy(socksProxy)
	}

	var resp *resty.Response
	resp, err = cli.R().
		Get(reqURL)
	if err != nil {
		logger.Sugar().Error("req err:", err)
		return nil, err
	}

	if !resp.IsSuccess() {
		logger.Sugar().Error("resp err: ", resp.StatusCode())
		return nil, err
	}

	queryParams, err := url.ParseQuery(resp.String())
	if err != nil {
		logger.Sugar().Error("decode param err: ", err)
		return nil, err
	}

	var token = make(map[string]interface{})
	for key, values := range queryParams {
		if len(values) > 0 {
			token[key] = values[0]
		}
	}

	return token, nil
}

func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	var err error

	socksProxy := os.Getenv("ENV_RECAPTCHA_REQUEST_PROXY")

	cli := resty.New()

	if socksProxy != "" {
		cli = cli.SetProxy(socksProxy)
	}

	var resp *resty.Response
	resp, err = cli.R().
		SetHeader("accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Get(constant.GithubGetUserInfoURL)
	if err != nil {
		logger.Sugar().Error("req err:", err)
		return nil, err
	}

	if !resp.IsSuccess() {
		logger.Sugar().Error("resp error: ", resp.StatusCode())
		return nil, err
	}

	var userInfo = make(map[string]interface{})
	if err = json.Unmarshal(resp.Body(), &userInfo); err != nil {
		logger.Sugar().Error("decode param err: ", err)
		return nil, err
	}

	return userInfo, nil
}
