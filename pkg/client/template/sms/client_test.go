package sms

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	valuedef "github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"

	"bou.ke/monkey"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/NpoolPlatform/third-middleware/pkg/testinit"

	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/sms"
	usedfor "github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	crud "github.com/NpoolPlatform/third-manager/pkg/crud/v1/template/sms"
	"github.com/stretchr/testify/assert"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

var data = mgrpb.SMSTemplate{
	ID:      uuid.NewString(),
	AppID:   uuid.NewString(),
	LangID:  uuid.NewString(),
	UsedFor: usedfor.UsedFor_KYCApproved,
	Subject: uuid.NewString(),
	Message: uuid.NewString(),
}

func getSMSTemplate(t *testing.T) {
	_, err := crud.Create(context.Background(), &mgrpb.SMSTemplateReq{
		ID:      &data.ID,
		AppID:   &data.AppID,
		LangID:  &data.LangID,
		UsedFor: &data.UsedFor,
		Subject: &data.Subject,
		Message: &data.Message,
	})
	assert.Nil(t, err)

	info, err := GetSMSTemplate(context.Background(), data.ID)
	if assert.Nil(t, err) {
		assert.Equal(t, &data, info)
	}
}

func getSMSTemplateOnly(t *testing.T) {
	info, err := GetSMSTemplateOnly(context.Background(), &mgrpb.Conds{
		ID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.ID,
		},
	})
	if assert.Nil(t, err) {
		assert.Equal(t, &data, info)
	}
}
func getSMSTemplates(t *testing.T) {
	infos, total, err := GetSMSTemplates(context.Background(), &mgrpb.Conds{
		ID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.ID,
		},
	}, 0, 1)
	if assert.Nil(t, err) {
		assert.Equal(t, total, uint32(1))
		assert.Equal(t, infos[0], &data)
	}
}

func TestClient(t *testing.T) {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}

	gport := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)

	monkey.Patch(grpc2.GetGRPCConn, func(service string, tags ...string) (*grpc.ClientConn, error) {
		return grpc.Dial(fmt.Sprintf("localhost:%v", gport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	t.Run("getSMSTemplate", getSMSTemplate)
	t.Run("getSMSTemplateOnly", getSMSTemplateOnly)
	t.Run("getSMSTemplates", getSMSTemplates)
}
