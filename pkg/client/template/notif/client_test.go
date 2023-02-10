package notif

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

	usedfor "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/notif"
	crud "github.com/NpoolPlatform/third-manager/pkg/crud/v1/template/notif"
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

var data = mgrpb.NotifTemplate{
	ID:      uuid.NewString(),
	AppID:   uuid.NewString(),
	LangID:  uuid.NewString(),
	UsedFor: usedfor.EventType_KYCApproved,
	Title:   uuid.NewString(),
	Content: uuid.NewString(),
}

func getNotifTemplate(t *testing.T) {
	_, err := crud.Create(context.Background(), &mgrpb.NotifTemplateReq{
		ID:      &data.ID,
		AppID:   &data.AppID,
		LangID:  &data.LangID,
		UsedFor: &data.UsedFor,
		Title:   &data.Title,
		Content: &data.Content,
	})
	assert.Nil(t, err)

	info, err := GetNotifTemplate(context.Background(), data.ID)
	if assert.Nil(t, err) {
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = info.UpdatedAt
		assert.Equal(t, &data, info)
	}
}

func getNotifTemplateOnly(t *testing.T) {
	info, err := GetNotifTemplateOnly(context.Background(), &mgrpb.Conds{
		ID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.ID,
		},
	})
	if assert.Nil(t, err) {
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = info.UpdatedAt
		assert.Equal(t, &data, info)
	}
}
func getNotifTemplates(t *testing.T) {
	infos, total, err := GetNotifTemplates(context.Background(), &mgrpb.Conds{
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
	t.Run("getNotifTemplate", getNotifTemplate)
	t.Run("getNotifTemplateOnly", getNotifTemplateOnly)
	t.Run("getNotifTemplates", getNotifTemplates)
}