package content

import "strings"

type Mark string

const (
	MarkSup Mark = "sup"
	MarkSub Mark = "sub"
)

type Inline struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Marks   []Mark `json:"marks,omitempty"`
	Latex   string `json:"latex,omitempty"`
	Display string `json:"display,omitempty"`
}

type Block struct {
	Type     string   `json:"type"`
	Children []Inline `json:"children,omitempty"`
}

type RichContent struct {
	Version int     `json:"version"`
	Blocks  []Block `json:"blocks,omitempty"`
}

func New() RichContent {
	return RichContent{Version: 1}
}

func (c RichContent) PlainText() string {
	if len(c.Blocks) == 0 {
		return ""
	}

	lines := make([]string, 0, len(c.Blocks))
	for _, block := range c.Blocks {
		if block.Type != "paragraph" {
			continue
		}
		var sb strings.Builder
		for _, child := range block.Children {
			switch child.Type {
			case "text":
				sb.WriteString(child.Text)
			case "math":
				sb.WriteString(child.Latex)
			}
		}
		lines = append(lines, strings.TrimSpace(sb.String()))
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func (c RichContent) Empty() bool {
	for _, block := range c.Blocks {
		for _, child := range block.Children {
			switch child.Type {
			case "text":
				if strings.TrimSpace(child.Text) != "" {
					return false
				}
			case "math":
				if strings.TrimSpace(child.Latex) != "" {
					return false
				}
			}
		}
	}
	return true
}
