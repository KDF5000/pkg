package larkbot

import (
	"fmt"
	"log"
	"testing"
)

func TestSendCardMessage(t *testing.T) {
	option := BotOption{
		AppID:       "xxx",
		AppSecret:   "xxx",
		VerifyToken: "xxx", // for event-subscribtion
		EncryptKey:  "xxx", // for event-subscribtion
	}

	bot := NewLarkBot(option)
	userInfo, err := bot.GetUserID("xxx@xxxx.com")
	if err != nil {
		t.Fatal(err)
	}
	// card, err := builder.BuildForm()
	builder := NewCardMessageBuilder()
	builder.Header("标题", "blue").
		Field("我是一个单独的div block").
		DIVBlock().
		Field("我是一个DivBlock里的field1").
		Field("我是一个DivBlock里的field2").
		Field(fmt.Sprintf("我要at<at email=xx@xx.com></at>")).
		HRBlock().
		Button("单独button", "", nil, "get").
		HRBlock().
		ButtonBlock().
		Button("button1", "", nil, "get").
		Button("button2", "", nil, "get")

	card, err := builder.Build()
	if err != nil {
		t.Fatal(err)
		return
	}

	err = bot.SendCardMessage(IDTypeUserID, userInfo.UserID, "", card)
	if err != nil {
		log.Fatal(err)
	}

}
