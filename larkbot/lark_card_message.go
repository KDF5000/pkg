package larkbot

import (
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type DivBlock struct {
    fields []protocol.FieldForm
}

type ButtonBlock struct {
    buttons []protocol.ActionElement
}

type CardMessageBuilder struct {
    builder *message.CardBuilder

    curButtonBlock *ButtonBlock
    curDivBlock *DivBlock
}

func NewCardMessageBuilder() *CardMessageBuilder {
    builder := &message.CardBuilder{}
    builder.SetConfig(protocol.ConfigForm{
        MinVersion: protocol.VersionForm{},
        WideScreenMode: true,
    })

    return &CardMessageBuilder{
        builder: builder,
    }
}

func (cb *CardMessageBuilder) Header(content, color string) *CardMessageBuilder {
    line := 1
    title := protocol.TextForm {
        Tag: protocol.PLAIN_TEXT_E,
        Content: &content,
        Lines: &line,
    }

    cb.builder.AddHeader(title, color)
    return cb
}

func (cb *CardMessageBuilder) sinkBlock()  {
    if cb.curDivBlock != nil && len(cb.curDivBlock.fields) > 0 {
        cb.builder.AddDIVBlock(nil, cb.curDivBlock.fields, nil)
        cb.curDivBlock = nil
    }

    if cb.curButtonBlock != nil && len(cb.curButtonBlock.buttons) > 0 {
	    cb.builder.AddActionBlock(cb.curButtonBlock.buttons)
        cb.curButtonBlock = nil
    }
}

func (cb *CardMessageBuilder) DIVBlock() *CardMessageBuilder {
    cb.sinkBlock()
    cb.curDivBlock = &DivBlock{}
    return cb
}

// content support markdown format
func (cb *CardMessageBuilder) Field(content string) *CardMessageBuilder {
    if cb.curDivBlock != nil {
        field := message.NewField(false, message.NewMDText(content, nil, nil, nil))
        cb.curDivBlock.fields = append(cb.curDivBlock.fields, *field)
        return cb
    }

    cb.builder.AddDIVBlock(nil, []protocol.FieldForm{
        *message.NewField(false, message.NewMDText(content, nil, nil, nil)),
    }, nil)

    return cb
}


func (cb *CardMessageBuilder) HRBlock() *CardMessageBuilder {
    cb.sinkBlock()
    cb.builder.AddHRBlock()
    return cb
}

func (cb *CardMessageBuilder) ButtonBlock() *CardMessageBuilder {
    cb.sinkBlock()
    cb.curButtonBlock = &ButtonBlock{}
    return cb
}

func (cb *CardMessageBuilder) Button(name, actionUrl string, params map[string]string, method string) *CardMessageBuilder {
    if cb.curButtonBlock != nil {
        button := message.NewButton(message.NewMDText(name, nil, nil, nil), &actionUrl, nil, params, protocol.PRIMARY, nil, method)
        cb.curButtonBlock.buttons = append(cb.curButtonBlock.buttons, button)
        return cb
    }

	cb.builder.AddActionBlock([]protocol.ActionElement{
        message.NewButton(message.NewMDText(name, nil, nil, nil), &actionUrl, nil, params, protocol.PRIMARY, nil, method),
	})

    return cb
}

func (cb *CardMessageBuilder) Build() (card *protocol.CardForm, err error){
    cb.sinkBlock()

    return cb.builder.BuildForm()
}
