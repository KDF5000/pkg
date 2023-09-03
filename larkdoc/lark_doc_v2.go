package larkdoc

import (
	"context"
	"fmt"
	"log"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

type DocV2Option struct {
	AppID     string
	AppSecret string

	// not necessary
	VerifyToken   string
	EncryptKey    string
	HelpDeskID    string
	HelpDeskToken string
}

type LarkDocV2 struct {
	option DocOption

	client *lark.Client
}

func NewLarkDocV2(opt DocOption) *LarkDocV2 {
	return &LarkDocV2{
		option: opt,
		client: lark.NewClient(opt.AppID, opt.AppSecret),
	}
}

func (d *LarkDocV2) GetRawContent(ctx context.Context, docToken string) (*DocRawContent, error) {
	buidler := larkdocx.NewRawContentDocumentReqBuilder()
	buidler.DocumentId(docToken)
	resp, err := d.client.Docx.Document.RawContent(ctx, buidler.Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("%s, requestid: %s", resp.Error(), resp.RequestId())
	}

	return &DocRawContent{
		Content: resp.Data.Content,
	}, nil
}

func (d *LarkDocV2) GetAllBlocks(ctx context.Context, docToken string) (*DocContentV2, error) {
	buidler := larkdocx.NewListDocumentBlockReqBuilder()
	buidler.DocumentId(docToken)
	hasMore := true
	var content DocContentV2
	for hasMore {
		resp, err := d.client.Docx.DocumentBlock.List(ctx, buidler.Build())
		if err != nil {
			return nil, err
		}

		if !resp.Success() {
			return nil, fmt.Errorf("%s, requestid: %s", resp.Error(), resp.RequestId())
		}

		content.Blocks = append(content.Blocks, resp.Data.Items...)
		hasMore = *resp.Data.HasMore
		if hasMore {
			buidler.PageToken(*resp.Data.PageToken)
		}
	}

	return &content, nil
}

func (d *LarkDocV2) BlocksIterator(ctx context.Context, dockToken string) (*larkdocx.ListDocumentBlockIterator, error) {
	builder := larkdocx.NewListDocumentBlockReqBuilder()
	builder.DocumentId(dockToken)

	log.Printf("builder: %+v", *(builder.Build()))

	return d.client.Docx.DocumentBlock.ListByIterator(ctx, builder.Build())
}

func (d *LarkDocV2) GetBasicInfo(ctx context.Context, docToken string) (*larkdocx.Document, error) {
	docReq := larkdocx.NewGetDocumentReqBuilder()
	docReq.DocumentId(docToken)

	resp, err := d.client.Docx.Document.Get(ctx, docReq.Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("%s, requestid: %s", resp.Error(), resp.RequestId())
	}

	return resp.Data.Document, nil
}

func (d *LarkDocV2) GetMeta(ctx context.Context, docToken string) (*DocMetaV2, error) {
	docReq := larkdrive.NewRequestDocBuilder()
	docReq.DocType("docx")
	docReq.DocToken(docToken)

	metaReq := larkdrive.NewMetaRequestBuilder()
	metaReq.RequestDocs([]*larkdrive.RequestDoc{docReq.Build()})
	metaReq.WithUrl(true)

	buidler := larkdrive.NewBatchQueryMetaReqBuilder()
	buidler.MetaRequest(metaReq.Build())
	resp, err := d.client.Drive.Meta.BatchQuery(ctx, buidler.Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("%s, requestid: %s", resp.Error(), resp.RequestId())
	}

	if len(resp.Data.Metas) <= 0 {
		if len(resp.Data.FailedList) > 0 {
			return nil, fmt.Errorf("failed to get doc meta, failed code: %d", *resp.Data.FailedList[0].Code)
		}
		return nil, fmt.Errorf("failed to get doc meta")
	}

	meta := resp.Data.Metas[0]
	if meta == nil {
		return nil, fmt.Errorf("failed to get doc meta")
	}

	return &DocMetaV2{meta}, nil
}

// insert a block at position pos
// block is a json string marshaled from Block struct
func (d *LarkDocV2) InsertBlock(ctx context.Context, docToken string, revision int, pos int, blocks []*larkdocx.Block) error {
	buidler := larkdocx.NewCreateDocumentBlockChildrenReqBodyBuilder()
	buidler.Index(pos)
	buidler.Children(blocks)

	req := larkdocx.NewCreateDocumentBlockChildrenReqBuilder()
	req.DocumentId(docToken)
	req.DocumentRevisionId(revision)
	// insert block into the page block tree
	req.BlockId(docToken)
	req.Body(buidler.Build())

	resp, err := d.client.Docx.DocumentBlockChildren.Create(ctx, req.Build())
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("%s, requestid: %s", resp.Error(), resp.RequestId())
	}

	return nil
}
