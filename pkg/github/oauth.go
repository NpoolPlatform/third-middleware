package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"
	constant "github.com/NpoolPlatform/third-middleware/pkg/const"
	"github.com/go-resty/resty/v2"
)

func GetAccessToken(clientID, clientSecret, code string) (*npool.AccessTokenInfo, error) {
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

	errStr := queryParams.Get("error_description")
	if errStr != "" {
		return nil, fmt.Errorf("%s", errStr)
	}

	return &npool.AccessTokenInfo{
		AccessToken:        queryParams.Get("access_token"),
		RefreshAccessToken: queryParams.Get("refresh_token"),
		Scope:              queryParams.Get("scope"),
		TokenType:          queryParams.Get("token_type"),
	}, nil
}

func GetUserInfo(accessToken string) (*npool.ThirdUserInfo, error) {
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
		return nil, fmt.Errorf("resp error: %v", resp.StatusCode())
	}

	var userInfo = make(map[string]interface{})
	if err = json.Unmarshal(resp.Body(), &userInfo); err != nil {
		logger.Sugar().Error("decode param err: ", err)
		return nil, err
	}

	if userInfo["message"] != nil {
		messageByte, err := json.Marshal(userInfo["message"])
		if err != nil {
			return nil, fmt.Errorf("get message error")
		}
		errStr := string(messageByte)
		return nil, fmt.Errorf("%s", errStr)
	}

	idByte, err := json.Marshal(userInfo["id"])
	if err != nil {
		return nil, fmt.Errorf("get id error")
	}
	nameByte, err := json.MarshalIndent(userInfo["name"], "", "  ")
	if err != nil {
		return nil, fmt.Errorf("get name error")
	}
	loginByte, err := json.MarshalIndent(userInfo["login"], "", "  ")
	if err != nil {
		return nil, fmt.Errorf("get login error")
	}
	avatarURLByte, err := json.MarshalIndent(userInfo["avatar_url"], "", "  ")
	if err != nil {
		return nil, fmt.Errorf("get avatarURL error")
	}
	info := &npool.ThirdUserInfo{
		ID:        string(idByte),
		Name:      string(nameByte),
		Login:     string(loginByte),
		AvatarURL: string(avatarURLByte),
	}

	return info, nil
}
