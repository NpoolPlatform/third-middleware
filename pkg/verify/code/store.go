package code

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
)

const (
	redisTimeout = 5 * time.Second
)

type UserCode struct {
	AppID       uuid.UUID
	Account     string
	AccountType string
	UsedFor     string
	Code        string
	NextAt      time.Time
	ExpireAt    time.Time
}

func (c *UserCode) Key() string {
	return fmt.Sprintf("verify-code:%v:%v:%v:%v", c.AppID, c.AccountType, c.Account, c.UsedFor)
}

func CreateCodeCache(ctx context.Context, code *UserCode) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	body, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("fail marshal code: %v", err)
	}

	err = cli.Set(ctx, code.Key(), body, time.Until(code.ExpireAt)).Err()
	if err != nil {
		return fmt.Errorf("fail create code cache: %v", err)
	}

	return nil
}

func VerifyCodeCache(ctx context.Context, code *UserCode) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, code.Key()).Result()
	if err == redis.Nil {
		return fmt.Errorf("code not found %v in redis", code.Key())
	} else if err != nil {
		return fmt.Errorf("fail get code: %v", err)
	}

	user := UserCode{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return fmt.Errorf("fail unmarshal val: %v", err)
	}

	if user.Code != code.Code {
		return fmt.Errorf("invalid code")
	}

	if time.Now().After(user.ExpireAt) {
		return fmt.Errorf("code expired")
	}

	err = DeleteCodeCache(ctx, code)
	if err != nil {
		return fmt.Errorf("fail delete code: %v", err)
	}

	return nil
}

func Nextable(ctx context.Context, code *UserCode) (bool, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return false, fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, code.Key()).Result()
	if err == redis.Nil {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("fail get code: %v", err)
	}

	user := UserCode{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return false, fmt.Errorf("fail unmarshal val: %v", err)
	}

	if !time.Now().After(user.NextAt) {
		return false, nil
	}

	return true, nil
}

func DeleteCodeCache(ctx context.Context, code *UserCode) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	err = cli.Del(ctx, code.Key()).Err()
	if err != nil {
		return fmt.Errorf("fail delete code: %v", err)
	}

	return nil
}
