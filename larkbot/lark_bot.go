package larkbot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type IDType protocol.UserType

const (
	IDTypeUnknown IDType = IDType(protocol.UserTypeUnknown)
	IDTypeOpenID  IDType = IDType(protocol.UserTypeOpenID)
	IDTypeUserID  IDType = IDType(protocol.UserTypeUserID)
	IDTypeEmail   IDType = IDType(protocol.UserTypeEmail)
	IDTypeChatID  IDType = IDType(protocol.UserTypeChatID)

	Email2IDPath protocol.OpenApiPath = "/open-apis/user/v4/email2id"
)

type Email2UIDReq struct {
	Email string `json:"email"`
}

type UserInfo struct {
	OpenID string `json:"open_id"`
	UserID string `json:"user_id"`
}

type Email2UIDResp struct {
	Code    int      `json:"code"`
	Message string   `json:"msg"`
	Data    UserInfo `json:"data"`
}

type BotOption struct {
	AppID       string
	AppSecret   string
	VerifyToken string
	EncryptKey  string
	TenantKey   string
}

type LarkBot struct {
	option BotOption
}

func NewLarkBot(opt BotOption) *LarkBot {
	larkBot := &LarkBot{option: opt}

	// init sdk
	larkBot.init()
	larkBot.init()

	return larkBot
}

func (bot *LarkBot) init() {
	appConf := appconfig.AppConfig{
		AppID:       bot.option.AppID,
		AppType:     protocol.InternalApp,
		AppSecret:   bot.option.AppSecret,
		VerifyToken: bot.option.VerifyToken,
		EncryptKey:  bot.option.EncryptKey,
	}

	appconfig.Init(appConf)
}

func (bot *LarkBot) SendTextMessage(idType IDType, chatID, rootID, msg string) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserType(idType),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := message.SendTextMessage(ctx, bot.option.TenantKey, bot.option.AppID, user, rootID, msg)
	if err != nil {
		return fmt.Errorf("send text failed, %v", err)
	}

	return nil
}

func (bot *LarkBot) SendImageMessage(idType IDType, chatID, rootID, imgUrl string) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserType(idType),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := message.SendImageMessage(ctx, bot.option.TenantKey, bot.option.AppID, user, rootID, imgUrl, "", "")
	if err != nil {
		return fmt.Errorf("send image failed, %v", err)
	}

	return nil
}

func (bot *LarkBot) SendRichTextMessage(idType IDType, chatID, rootID string, forms map[protocol.Language]*protocol.RichTextForm) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserType(idType),
	}
	//send rich text
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err := message.SendRichTextMessage(ctx, bot.option.TenantKey, bot.option.AppID, user, rootID, forms)
	if err != nil {
		return fmt.Errorf("send rich text failed[%v]", err)
	}

	return nil
}

func (bot *LarkBot) SendCardMessage(idType IDType, chatID, rootID string, card *protocol.CardForm) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserType(idType),
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err := message.SendCardMessage(ctx, bot.option.TenantKey, bot.option.AppID, user, rootID, *card, false)
	if err != nil {
		return fmt.Errorf("send card message failed[%v]", err)
	}

	return nil
}

func (bot *LarkBot) GetUserID(email string) (UserInfo, error) {
	req := Email2UIDReq{
		Email: email,
	}

	tenant_token, err := auth.GetTenantAccessToken(context.Background(), "", bot.option.AppID)
	if err != nil {
		return UserInfo{}, fmt.Errorf("GetTenantAccessToken error, %v", err)
	}

	headers := common.NewHeaderToken(tenant_token)
	data, code, err := common.DoHttpPostOApi(Email2IDPath, headers, req)
	if err != nil {
		return UserInfo{}, fmt.Errorf("DoHttpPostOApi error, code: %d, error: %v", code, err)
	}

	var resp Email2UIDResp
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return UserInfo{}, fmt.Errorf("GetUserID Unmarshal error, %s", err)
	}

	if resp.Code != 0 {
		return UserInfo{}, fmt.Errorf("get UserID error, %+v", resp.Message)
	}

	return resp.Data, nil
}
