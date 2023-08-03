package oauth

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"
	github "github.com/NpoolPlatform/third-middleware/pkg/github"
)

func (h *Handler) GetAccessToken(ctx context.Context) (*npool.AccessTokenInfo, error) {
	if h.Code == "" {
		return nil, fmt.Errorf("invalid code")
	}

	switch h.ClientName {
	case basetypes.SignMethod_Github:
		info, err := github.GetAccessToken(h.ClientID, h.ClientSecret, h.Code)
		if err != nil {
			return nil, err
		}
		return info, nil
	default:
		return nil, fmt.Errorf("unsupport oauth")
	}
}

func (h *Handler) GetThirdUserInfo(ctx context.Context) (*npool.ThirdUserInfo, error) {
	if h.AccessToken == "" {
		return nil, fmt.Errorf("invalid accesstoken")
	}

	switch h.ClientName {
	case basetypes.SignMethod_Github:
		info, err := github.GetUserInfo(h.AccessToken)
		if err != nil {
			return nil, err
		}
		return info, nil
	default:
		return nil, fmt.Errorf("unsupport oauth")
	}
}
