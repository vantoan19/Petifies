package valueobjects

type TextContent struct {
	content string
}

func NewTextContent(content string) TextContent {
	return TextContent{content: content}
}

func (t TextContent) IsEmpty() bool {
	return t.content == ""
}

func (t TextContent) Content() string {
	return t.content
}
