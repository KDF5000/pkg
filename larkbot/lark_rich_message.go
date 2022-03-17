package larkbot

import (
	"log"

	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

var langMap = map[string]protocol.Language{
	"cn": protocol.ZhCN,
	"en": protocol.EnUS,
	"jp": protocol.JaJP,
}

type RichMessageBuilder struct {
	forms map[protocol.Language]*protocol.RichTextForm

	curContent *protocol.RichTextContent
	// element`s for last/curent line
	// Used to build a new line when start next line or
	// try to build final form
	elements []*protocol.RichTextElementForm
}

func NewMessageBuilder() *RichMessageBuilder {
	builder := &RichMessageBuilder{
		forms: make(map[protocol.Language]*protocol.RichTextForm),
	}
	return builder
}

func (builder *RichMessageBuilder) sinkElements() {
	if builder.curContent != nil {
		if len(builder.elements) > 0 {
			builder.curContent.AddElementBlock(builder.elements...)
		}
	}

	builder.elements = builder.elements[:0]
}

func (builder *RichMessageBuilder) AddForm(lan, title string) *RichMessageBuilder {
	if builder.forms == nil {
		builder.forms = make(map[protocol.Language]*protocol.RichTextForm)
	}

	language, ok := langMap[lan]
	if !ok {
		// do not panic ?
		log.Printf("`%s` not supported, any operation for this form will be not worked!", lan)
		// panic(fmt.Sprintf("language %s is not supported", lan))
		return builder
	}

	builder.sinkElements()
	content := message.NewRichTextContent()
	builder.forms[language] = message.NewRichTextForm(&title, content)
	builder.curContent = content
	return builder
}

func (builder *RichMessageBuilder) NewLine() *RichMessageBuilder {
	builder.sinkElements()

	return builder
}

func (builder *RichMessageBuilder) BlankLine() *RichMessageBuilder {
	builder.sinkElements()
	if builder.curContent != nil {
		builder.curContent.AddElementBlock()
	}

	return builder
}

func (builder *RichMessageBuilder) TextElement(content string) *RichMessageBuilder {
	tag := message.NewTextTag(content, true, 1)
	builder.elements = append(builder.elements, tag)
	return builder
}

func (builder *RichMessageBuilder) LinkElement(text, url string) *RichMessageBuilder {
	tag := message.NewATag(text, true, url)
	builder.elements = append(builder.elements, tag)
	return builder
}

func (builder *RichMessageBuilder) AtElement(username, uid string) *RichMessageBuilder {
	tag := message.NewAtTag(username, uid)
	builder.elements = append(builder.elements, tag)
	return builder
}

func (builder *RichMessageBuilder) Build() map[protocol.Language]*protocol.RichTextForm {
	builder.sinkElements()

	builder.curContent = nil
	return builder.forms
}
