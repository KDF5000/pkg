package larkbot

import (
	"context"
	"testing"

	"github.com/larksuite/botframework-go/SDK/auth"
)

func TestToken(t *testing.T) {
	option := BotOption{
		AppID:       "xxx",
		AppSecret:   "xxx",
		VerifyToken: "xxx", // for event-subscribtion
		EncryptKey:  "xxx", // for event-subscribtion
	}

	_ = NewLarkBot(option)
	app_access_token, err := auth.GetAppAccessToken(context.TODO(), option.AppID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("app_access_token: %s", app_access_token)

	tenant_key, err := auth.GetTenantAccessToken(context.Background(), "", option.AppID)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("tenant_key: %s", tenant_key)
}

func TestGetUserID(t *testing.T) {
	option := BotOption{
		AppID:       "xxx",
		AppSecret:   "xxx",
		VerifyToken: "xxx", // for event-subscribtion
		EncryptKey:  "xxx", // for event-subscribtion
	}

	larkbot := NewLarkBot(option)
	userInfo, err := larkbot.GetUserID("xx@xxx.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("open_id: %s, user_id: %s", userInfo.OpenID, userInfo.UserID)
}
