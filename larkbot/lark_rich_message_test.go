package larkbot

import (
	"log"
	"testing"
)

func TestSendMessage(t *testing.T) {
	option := BotOption{
		AppID:       "xxx",
		AppSecret:   "xxx",
		VerifyToken: "xxx", // for event-subscribtion
		EncryptKey:  "xxx", // for event-subscribtion
	}

	bot := NewLarkBot(option)
	userInfo, err := bot.GetUserID("xxx@xxx.com")
	if err != nil {
		t.Fatal(err)
	}

	builder := NewMessageBuilder()
	builder.AddForm("en", "xxx open a merge request").
		TextElement("Hey ").AtElement("xxx", userInfo.UserID).TextElement(", this MR is assigned to you.").
		BlankLine().
		NewLine().TextElement("Repository: xxx").
		NewLine().TextElement("Commit: Support unit index in replicated chunks....").
		NewLine().TextElement("State: opened").
		BlankLine().
		NewLine().TextElement("MergeRequest: ").LinkElement("https://baidu.com", "https://baidu.com").
		NewLine().TextElement("CommitUrl: ").LinkElement("https://baidu.com")

	err = bot.SendRichTextMessage(IDTypeUserID, userInfo.UserID, "", builder.Build())
	if err != nil {
		log.Fatal(err)
	}
}
