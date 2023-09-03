package larkdoc

import (
	doc "github.com/larksuite/oapi-sdk-go/service/doc/v2"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

// 1	页面 Block
// 2	文本 Block
// 3	标题 1 Block
// 4	标题 2 Block
// 5	标题 3 Block
// 6	标题 4 Block
// 7	标题 5 Block
// 8	标题 6 Block
// 9	标题 7 Block
// 10	标题 8 Block
// 11	标题 9 Block
// 12	无序列表 Block
// 13	有序列表 Block
// 14	代码块 Block
// 15	引用 Block
// 17	待办事项 Block
// 18	多维表格 Block
// 19	高亮块 Block
// 20	会话卡片 Block
// 21	流程图 & UML Block
// 22	分割线 Block
// 23	文件 Block
// 24	分栏 Block
// 25	分栏列 Block
// 26	内嵌 Block Block
// 27	图片 Block
// 28	开放平台小组件 Block
// 29	思维笔记 Block
// 30	电子表格 Block
// 31	表格 Block
// 32	表格单元格 Block
// 33	视图 Block
// 34	引用容器 Block
// 35	任务 Block
// 36	OKR Block
// 37	OKR Objective Block
// 38	OKR Key Result Block
// 39	OKR Progress Block
// 40	新版文档小组件 Block
// 41	Jira 问题 Block
// 42	Wiki 子目录 Block
// 999	未支持 Block
type DocxBlockType int

const (
	DocxBlockTypePage = iota + 1
	DocxBlockTypeText
	DocxBlockTypeHeading1
	DocxBlockTypeHeading2
	DocxBlockTypeHeading3
	DocxBlockTypeHeading4
	DocxBlockTypeHeading5
	DocxBlockTypeHeading6
	DocxBlockTypeHeading7
	DocxBlockTypeHeading8
	DocxBlockTypeHeading9
	DocxBlockTypeBullet
	DocxBlockTypeOrdered
	DocxBlockTypeCode
	DocxBlockTypeQuote
	DocxBlockTypeTodo
	DocxBlockTypeBitable
	DocxBlockTypeCallout
	DocxBlockTypeDiagram
	DocxBlockTypeDivider
	DocxBlockTypeFile
	DocxBlockTypeGrid
	DocxBlockTypeGridColumn
	DocxBlockTypeIframe
	DocxBlockTypeImag
	DocxBlockTypeISV
	DocxBlockTypeMindnote
	DocxBlockTypeSheet
	DocxBlockTypeTable
	DocxBlockTypeTableCell
	DocxBlockTypeView
	DocxBlockTypeQuoteContainer
	DocxBlockTypeTask
	DocxBlockTypeOKR
	DocxBlockTypeObjective
	DocxBlockTypeKeyResult
	DocxBlockTypeProgress
	DocxBlockTypeNewDocBlock
	DocxBlockTypeJira
	DocxBlockTypeWikiSubDir

	DocxBlockTypeUndefined = 999
)

type BlockType string
type ParagraphElementType string
type ListType string
type OperationRequestType string

const (
	BLOCK_TYPE_PARAGRAPH = "paragraph"
	BLOCK_TYPE_CALLOUT   = "callout"

	PARAGRAPH_ELEMENT_TYPE_TEXTRUN  = "textRun"
	PARAGRAPH_ELEMENT_TYPE_DOCSLINK = "docsLink"

	// number|bullet|checkBox|checkedBox|code
	LIST_TYPE_NUMBER     = "number"
	LIST_TYPE_BULLET     = "bullet"
	LIST_TYPE_CHECKBOX   = "checkBox"
	LIST_TYPE_CHECKEDBOX = "checkedBox"
	LIST_TYPE_CODE       = "code"

	// batchupdate
	OP_REQUEST_TYPE_INSERTBLOCK = "InsertBlocksRequestType"
)

type Location struct {
	ZoneId     string `json:"zoneId,omitempty"`
	StartIndex uint64 `json:"startIndex,omitempty"`
	EndIndex   uint64 `json:"endIndex,omitempty"`
}

type RGBColor struct {
	Red   int `json:"red,omitempty"`
	Green int `json:"green,omitempty"`
	Blue  int `json:"blue,omitempty"`
	Alpha int `json:"alpha,omitempty"`
}

type CalloutBlock struct {
	Location               Location `json:"location,omitempty"`
	CalloutEmojiId         string   `json:"calloutEmojiId,omitempty"`
	CalloutBackgroundColor RGBColor `json:"calloutBackgroundColor,omitempty"`
	CalloutBorderColor     RGBColor `json:"calloutBorderColor,omitempty"`
	ZoneId                 string   `json:"zoneId,omitempty"`
	Body                   struct {
		Blocks []Block `json:"blocks,omitempty"`
	} `json:"body,omitempty"`
}

type List struct {
	Type        ListType `json:"type" comments:"number|bullet|checkBox|checkedBox|code" `
	IndentLevel int      `json:"indentLevel,omitempty"`
	Number      int      `json:"number,omitempty"`
}

type ParagraphStyle struct {
	HeadingLevel int    `json:"headingLevel,omitempty"`
	Collapse     bool   `json:"collapse,omitempty"`
	List         List   `json:"list,omitempty"`
	Quote        bool   `json:"quote,omitempty"`
	Align        string `json:"align,omitempty"`
}

type Link struct {
	URL string `json:"url,omitempty"`
}

type TextStyle struct {
	Bold          bool     `json:"bold,omitempty"`
	Italic        bool     `json:"italic,omitempty"`
	StrikeThrough bool     `json:"strikeThrough,omitempty"`
	UnderLine     bool     `json:"underLine,omitempty"`
	CodeInline    bool     `json:"codeInline,omitempty"`
	BackColor     RGBColor `json:"backColor,omitempty"`
	TextColor     RGBColor `json:"textColor,omitempty"`
	Link          Link     `json:"link,omitempty"`
}

type TextRun struct {
	Text     string    `json:"text,omitempty"`
	Style    TextStyle `json:"style,omitempty"`
	LineId   string    `json:"lineId,omitempty"`
	Location Location  `json:"location,omitempty"`
}

type DocsLink struct {
	URL      string   `json:"url,omitempty"`
	Location Location `json:"location,omitempty"`
}

type ParagraphElement struct {
	// "textRun": {object(TextRun)},
	// "docsLink": {object(DocsLink)},
	// "person": {object(Person)},
	// "equation": {object(Equation)},
	// "reminder": {object(Reminder)},
	// "file": {object(File)},
	// "jira": {object(Jira)},
	// "undefinedElement": {object(UndefinedElement)}
	Type ParagraphElementType `json:"type,omitempty"`

	TextRun  *TextRun  `json:"textRun,omitempty"`
	DocsLink *DocsLink `json:"docsLink,omitempty"`
}

type ParagraphBlock struct {
	Location Location           `json:"location,omitempty"`
	Style    ParagraphStyle     `json:"style,omitempty"`
	Elements []ParagraphElement `json:"elements,omitempty"`
}

type Block struct {
	// paragraph 文本段落
	// gallery 图片
	// file 文件上传
	// chatGroup 群名片
	// table 格式化表格
	// horizontalLine 水平分割线
	// embeddedPage 内嵌网页，例如西瓜视频等
	// sheet 电子表格
	// bitable 数据表或看板
	// diagram 绘图或uml图
	// jira jira filter或jira issue
	// poll 投票
	// code 新代码块
	// docsApp 团队互动应用
	// callout 高亮块
	// undefinedBlock 未支持的block，全部用undefineBlock表示
	Type BlockType `json:"type,omitempty"`

	// one of the following blocks will be populated
	Paragraph *ParagraphBlock `json:"paragraph,omitempty"`
	Callout   *CalloutBlock   `json:"callout,omitempty"`
}

type DocTitle struct {
	Elements []ParagraphElement `json:"elements"`
	Location Location           `json:"location"`
	LineId   string             `json:"lineId"`
}

type DocBody struct {
	Blocks []Block `json:"blocks"`
}

type DocContent struct {
	Revision int      `json:"revision,omitempty"`
	Title    DocTitle `json:"title"`
	Body     DocBody  `json:"body"`
}

type DocContentV2 struct {
	Blocks []*larkdocx.Block `json:"blocks,omitempty"`
}

type DocMeta struct {
	*doc.DocMetaResult
}

type DocMetaV2 struct {
	*larkdrive.Meta
}

type DocRawContent struct {
	Content *string `json:"content"`
}

type InsertBlocksRequest struct {
	Payload  string         `json:"payload"`
	Location InsertLocation `json:"location"`
}

type InsertLocation struct {
	ZoneID      string `json:"zoneId"`
	Index       uint64 `json:"index"`
	StartOfZone bool   `json:"startOfZone"`
	EndOfZone   bool   `json:"endOfZone"`
}

type OperationRequest struct {
	RequestType OperationRequestType `json:"requestType"`

	InsertBlocksRequest *InsertBlocksRequest `json:"insertBlocksRequest,omitempty"`
}
