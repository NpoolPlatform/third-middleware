package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
)

func GetAccessToken(clientID, clientSecret, code string) (map[string]interface{}, error) {
	var req *http.Request
	var err error
	var getTokenURL = "https://github.com/login/oauth/access_token"
	reqBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s",
		clientID, clientSecret, code,
	)
	if req, err = http.NewRequest(http.MethodGet, getTokenURL+"?"+reqBody, nil); err != nil {
		logger.Sugar().Error("build req err: ", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		logger.Sugar().Error("send req err: ", err)
		return nil, err
	}

	var token = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		logger.Sugar().Error("resp err: ", err)
		return nil, err
	}
	return token, nil
}

func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	var req *http.Request
	var err error
	var getUserInfoURL = "https://api.github.com/user"
	if req, err = http.NewRequest(http.MethodGet, getUserInfoURL, nil); err != nil {
		logger.Sugar().Error("build req err: ", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		logger.Sugar().Error("send req err: ", err)
		return nil, err
	}

	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		logger.Sugar().Error("resp err: ", err)
		return nil, err
	}
	return userInfo, nil
}
