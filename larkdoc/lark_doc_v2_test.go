package larkdoc

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

func TestLarkDocV2(t *testing.T) {
	doc := NewLarkDocV2(DocOption{
		AppID:     "cli_a10d258394b8500b",
		AppSecret: "FRKdRyJnOIAHShoImOI1odEcw4m1Di40",
	})

	content, err := doc.GetRawContent(context.Background(), "CZJidAkBUonzOTx8cGtcyB6MnOd")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("RawContent: %s", *content.Content)

	docContent, err := doc.GetAllBlocks(context.Background(), "CZJidAkBUonzOTx8cGtcyB6MnOd")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("docContent: %+v", *docContent)
	for i := range docContent.Blocks {
		block := docContent.Blocks[i]
		if *(block.BlockType) == 1 {
			t.Logf("%+v", *block.Page)
		}
	}

	meta, err := doc.GetMeta(context.Background(), "CZJidAkBUonzOTx8cGtcyB6MnOd")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Meta: %+v", *(meta.Meta))
}

func TestLarkDocV2InsertBlock(t *testing.T) {
	doc := NewLarkDocV2(DocOption{
		AppID:     "cli_a10d258394b8500b",
		AppSecret: "FRKdRyJnOIAHShoImOI1odEcw4m1Di40",
	})

	basicInfo, err := doc.GetBasicInfo(context.TODO(), "CZJidAkBUonzOTx8cGtcyB6MnOd")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("baseicInfo: %s, %d, %s", *basicInfo.DocumentId, *basicInfo.RevisionId, *basicInfo.Title)

	dateText := larkdocx.NewTextRunBuilder().Content("2022-02-06").Build()
	text := larkdocx.NewTextBuilder().Elements([]*larkdocx.TextElement{
		larkdocx.NewTextElementBuilder().TextRun(dateText).Build(),
	}).Build()

	dateBlock := larkdocx.NewBlockBuilder()
	dateBlock.BlockType(DocxBlockTypeHeading3)
	dateBlock.Heading3(text)

	date, _ := json.Marshal(dateBlock.Build())
	log.Printf("%+v", string(date))

	err = doc.InsertBlock(context.Background(), "CZJidAkBUonzOTx8cGtcyB6MnOd", -1,
		0, []*larkdocx.Block{dateBlock.Build()})
	if err != nil {
		log.Fatal(err)
	}

	memoText := larkdocx.NewTextRunBuilder().
		Content("产品化、服务化思维对RD也是非常重要的能力，能够把技术能力通过产品的形式为用户提供服务比技术本身要更加难")

	memo := larkdocx.NewTextBuilder().Elements([]*larkdocx.TextElement{
		larkdocx.NewTextElementBuilder().TextRun(memoText.Build()).Build(),
	}).Build()

	memoBlock := larkdocx.NewBlockBuilder()
	memoBlock.BlockType(DocxBlockTypeBullet)
	memoBlock.Bullet(memo)

	err = doc.InsertBlock(context.Background(), "CZJidAkBUonzOTx8cGtcyB6MnOd", -1,
		0, []*larkdocx.Block{memoBlock.Build()})
	if err != nil {
		log.Fatal(err)
	}
}

func TestIterator(t *testing.T) {
	doc := NewLarkDocV2(DocOption{
		AppID:     "cli_a10d258394b8500b",
		AppSecret: "FRKdRyJnOIAHShoImOI1odEcw4m1Di40",
	})

	iter, err := doc.BlocksIterator(context.TODO(), "CZJidAkBUonzOTx8cGtcyB6MnOd")
	if err != nil {
		t.Fatal(err)
	}

	for {
		hasMore, block, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}

		if !hasMore {
			break
		}

		t.Logf("%+v", *block.BlockType)
	}
}
