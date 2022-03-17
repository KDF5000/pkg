package larkdoc

import (
	"context"
	"testing"
)

func TestLarkDoc(t *testing.T) {
	doc := NewLarkDoc(DocOption{
		AppID:     "xxx",
		AppSecret: "xxx",
	})

	content, err := doc.GetRawContent(context.Background(), "xxx")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("RawContent: %s", *content.Content)

	docContent, err := doc.GetContent(context.Background(), "xxx")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("docContent: %+v", docContent)
	for i := range docContent.Body.Blocks {
		block := &docContent.Body.Blocks[i]
		if block.Type == "callout" {
			t.Logf("Type: %s, Callout: %+v", block.Type, block.Callout)
		} else if block.Type == "paragraph" {
			t.Logf("Type: %s, Paragraph: %+v", block.Type, block.Paragraph)
		} else {
			t.Logf("unknown type %s", block.Type)
		}
	}

	meta, err := doc.GetMeta(context.Background(), "xxxx")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Meta: %+v", meta)
}

func TestLarkDocInsertBlock(t *testing.T) {

	doc := NewLarkDoc(DocOption{
		AppID:     "xxx",
		AppSecret: "xxxxx",
	})

	oldContent, err := doc.GetContent(context.Background(), "xxxx")
	if err != nil {
		t.Fatal(err)
	}

	var firstIndex uint64
	for i := range oldContent.Body.Blocks {
		block := &oldContent.Body.Blocks[i]
		if block.Type == BLOCK_TYPE_PARAGRAPH {
			firstIndex = block.Paragraph.Location.StartIndex
			break
		}
	}

	var dateBlock Block
	dateBlock.Type = BLOCK_TYPE_PARAGRAPH
	var paragraph ParagraphBlock
	paragraph.Style.HeadingLevel = 3
	paragraph.Elements = append(paragraph.Elements,
		ParagraphElement{Type: PARAGRAPH_ELEMENT_TYPE_TEXTRUN,
			TextRun: &TextRun{Text: "2022-02-06"}})
	dateBlock.Paragraph = &paragraph

	var memo Block
	memo.Type = BLOCK_TYPE_PARAGRAPH
	memo.Paragraph = &ParagraphBlock{}
	memo.Paragraph.Style.List.Type = LIST_TYPE_BULLET
	memo.Paragraph.Style.List.IndentLevel = 1
	memo.Paragraph.Elements = append(memo.Paragraph.Elements,
		ParagraphElement{Type: PARAGRAPH_ELEMENT_TYPE_TEXTRUN,
			TextRun: &TextRun{Text: "产品化、服务化思维对RD也是非常重要的能力，能够把技术能力通过产品的形式为用户提供服务比技术本身要更加难"}})

	err = doc.InsertBlock(context.Background(), "xxxx",
		oldContent.Revision, firstIndex, dateBlock, memo, memo)
	if err != nil {
		t.Fatal(err)
	}
}
