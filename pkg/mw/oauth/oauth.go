package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"
	github "github.com/NpoolPlatform/third-middleware/pkg/github"
)

func (h *Handler) GetAccessToken(ctx context.Context) (string, error) {
	if h.Code == "" {
		return "", fmt.Errorf("invalid code")
	}

	switch h.ClientName {
	case basetypes.SignMethod_Github:
		jsonMap, err := github.GetAccessToken(h.ClientID, h.ClientSecret, h.Code)
		if err != nil {
			return "", err
		}
		if jsonMap["error_description"] != nil {
			errByte, err := json.Marshal(jsonMap["error_description"])
			if err != nil {
				return "", fmt.Errorf("get error_description error")
			}
			return "", fmt.Errorf(string(errByte))
		}
		accessTokenByte, err := json.Marshal(jsonMap["access_token"])
		if err != nil {
			return "", fmt.Errorf("get accesstoken error")
		}

		return string(accessTokenByte), nil
	default:
		return "", fmt.Errorf("unsupport oauth")
	}
}

func (h *Handler) GetThirdUserInfo(ctx context.Context) (*npool.ThirdUserInfo, error) {
	if h.AccessToken == "" {
		return nil, fmt.Errorf("invalid accesstoken")
	}

	switch h.ClientName {
	case basetypes.SignMethod_Github:
		jsonMap, err := github.GetUserInfo(h.AccessToken)
		if err != nil {
			return nil, err
		}
		if jsonMap["message"] != nil {
			messageByte, err := json.Marshal(jsonMap["message"])
			if err != nil {
				return nil, fmt.Errorf("get message error")
			}
			return nil, fmt.Errorf(string(messageByte))
		}
		idByte, err := json.Marshal(jsonMap["id"])
		if err != nil {
			return nil, fmt.Errorf("get id error")
		}
		nameByte, err := json.Marshal(jsonMap["name"])
		if err != nil {
			return nil, fmt.Errorf("get name error")
		}
		loginByte, err := json.Marshal(jsonMap["login"])
		if err != nil {
			return nil, fmt.Errorf("get login error")
		}
		avatarURLByte, err := json.Marshal(jsonMap["avatar_url"])
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
	default:
		return nil, fmt.Errorf("unsupport oauth")
	}
}
