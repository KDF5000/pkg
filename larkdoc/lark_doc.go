package larkdoc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/larksuite/oapi-sdk-go/core"
	doc "github.com/larksuite/oapi-sdk-go/service/doc/v2"
)

type DocOption struct {
	AppID     string
	AppSecret string

	// not necessary
	VerifyToken   string
	EncryptKey    string
	HelpDeskID    string
	HelpDeskToken string
}

type LarkDoc struct {
	option DocOption

	// inherit doc service`s method
	*doc.Service
}

func NewLarkDoc(opt DocOption) *LarkDoc {
	larkDoc := &LarkDoc{
		option: opt,
	}
	appSettings := core.NewInternalAppSettings(
		core.SetAppCredentials(opt.AppID, opt.AppSecret),               // 必需
		core.SetAppEventKey(opt.VerifyToken, opt.EncryptKey),           // 非必需，订阅事件、消息卡片时必需
		core.SetHelpDeskCredentials(opt.HelpDeskID, opt.HelpDeskToken)) // 非必需，使用服务台API时必需
	larkDoc.Service = doc.NewService(core.NewConfig(core.DomainFeiShu, appSettings,
		core.SetLoggerLevel(core.LoggerLevelError)))
	return larkDoc
}

func (d *LarkDoc) GetRawContent(ctx context.Context, docToken string) (*DocRawContent, error) {
	ctxW := core.WrapContext(ctx)
	reqCall := d.Docs.RawContent(ctxW)
	reqCall.SetDocToken(docToken)
	res, err := reqCall.Do()
	if err != nil {
		return nil, err
	}

	return &DocRawContent{
		Content: &res.Content,
	}, nil
}

func (d *LarkDoc) GetContent(ctx context.Context, docToken string) (*DocContent, error) {
	ctxW := core.WrapContext(ctx)
	reqCall := d.Docs.Content(ctxW)
	reqCall.SetDocToken(docToken)
	res, err := reqCall.Do()
	if err != nil {
		return nil, err
	}

	docContent := &DocContent{
		Revision: res.Revision,
	}
	docContent.Revision = res.Revision
	if err := json.Unmarshal([]byte(res.Content), docContent); err != nil {
		return nil, err
	}

	return docContent, nil
}

func (d *LarkDoc) GetMeta(ctx context.Context, docToken string) (*DocMeta, error) {
	ctxW := core.WrapContext(ctx)
	metaCall := d.Docs.Meta(ctxW)
	metaCall.SetDocToken(docToken)
	meta, err := metaCall.Do()
	if err != nil {
		return nil, err
	}

	return &DocMeta{meta}, nil
}

// insert a block at position pos
// block is a json string marshaled from Block struct
func (d *LarkDoc) InsertBlock(ctx context.Context, docToken string, revision int, pos uint64, block ...Block) error {
	var db DocBody
	db.Blocks = append(db.Blocks, block...)
	payload, err := json.Marshal(&db)
	if err != nil {
		return err
	}

	var req InsertBlocksRequest
	req.Payload = string(payload)
	req.Location = InsertLocation{
		ZoneID: "0",
		Index:  pos,
	}

	var opRequest OperationRequest
	opRequest.RequestType = OP_REQUEST_TYPE_INSERTBLOCK
	opRequest.InsertBlocksRequest = &req
	data, err := json.Marshal(&opRequest)
	if err != nil {
		return err
	}

	log.Printf("request: %s", string(data))

	// send a batch update request
	body := &doc.DocBatchUpdateReqBody{Revision: revision}
	body.Requests = append(body.Requests, string(data))
	ctxW := core.WrapContext(ctx)
	reqCall := d.Docs.BatchUpdate(ctxW, body)
	reqCall.SetDocToken(docToken)
	_, err = reqCall.Do()
	return err
}
