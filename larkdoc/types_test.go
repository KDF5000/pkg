package larkdoc

import (
	"encoding/json"
	"testing"
)

func TestParagraph(t *testing.T) {
	var date Block
	date.Type = "paragraph"
	date.Paragraph = &ParagraphBlock{}
	date.Paragraph.Style.HeadingLevel = 3
	date.Paragraph.Elements = append(date.Paragraph.Elements, ParagraphElement{
		Type: "textRun",
		TextRun: &TextRun{
			Text: "2022-2-5",
		},
	})

	data, err := json.Marshal(&date)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", data)
}
