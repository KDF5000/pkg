package larkdoc

import (
	doc "github.com/larksuite/oapi-sdk-go/service/doc/v2"
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

type DocMeta struct {
	*doc.DocMetaResult
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
