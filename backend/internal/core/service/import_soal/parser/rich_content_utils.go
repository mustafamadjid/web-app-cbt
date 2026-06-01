package parser

import (
	"strings"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
)

func appendContent(target *content.RichContent, block content.RichContent) {
	if target.Version == 0 {
		target.Version = 1
	}
	target.Blocks = append(target.Blocks, block.Blocks...)
}

func cloneMarks(marks []content.Mark) []content.Mark {
	if len(marks) == 0 {
		return nil
	}
	out := make([]content.Mark, len(marks))
	copy(out, marks)
	return out
}

func mergeMarks(base, extra []content.Mark) []content.Mark {
	out := cloneMarks(base)
	for _, mark := range extra {
		exists := false
		for _, current := range out {
			if current == mark {
				exists = true
				break
			}
		}
		if !exists {
			out = append(out, mark)
		}
	}
	return out
}

func marksFromRun(node xmlNode) []content.Mark {
	props := findFirstNode(node, "rPr")
	if props == nil {
		return nil
	}
	vertAlign := findFirstNode(*props, "vertAlign")
	if vertAlign == nil {
		return nil
	}
	switch attrValue(*vertAlign, "val") {
	case "superscript":
		return []content.Mark{content.MarkSup}
	case "subscript":
		return []content.Mark{content.MarkSub}
	default:
		return nil
	}
}

func leadingImageFilename(raw string) (string, bool) {
	fields := strings.Fields(strings.TrimSpace(raw))
	if len(fields) == 0 {
		return "", false
	}
	name := fields[0]
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".png"),
		strings.HasSuffix(lower, ".jpg"),
		strings.HasSuffix(lower, ".jpeg"),
		strings.HasSuffix(lower, ".gif"),
		strings.HasSuffix(lower, ".webp"):
		return name, true
	default:
		return "", false
	}
}

func trimContentPrefix(value content.RichContent, prefixLen int) content.RichContent {
	if prefixLen <= 0 || len(value.Blocks) == 0 {
		return value
	}

	remaining := prefixLen
	result := content.New()
	for _, block := range value.Blocks {
		nextBlock := content.Block{Type: block.Type}
		for _, child := range block.Children {
			switch child.Type {
			case "text":
				if remaining > 0 {
					runes := []rune(child.Text)
					if len(runes) <= remaining {
						remaining -= len(runes)
						continue
					}
					child.Text = string(runes[remaining:])
					remaining = 0
				}
				if child.Text != "" {
					nextBlock.Children = append(nextBlock.Children, child)
				}
			default:
				nextBlock.Children = append(nextBlock.Children, child)
			}
		}
		if len(nextBlock.Children) > 0 {
			result.Blocks = append(result.Blocks, nextBlock)
		}
	}

	return trimLeadingWhitespaceContent(result)
}

func trimLeadingWhitespaceContent(value content.RichContent) content.RichContent {
	for blockIdx := range value.Blocks {
		for childIdx := range value.Blocks[blockIdx].Children {
			child := &value.Blocks[blockIdx].Children[childIdx]
			if child.Type != "text" {
				return value
			}
			trimmed := strings.TrimLeft(child.Text, " \t")
			if trimmed == "" {
				child.Text = ""
				continue
			}
			child.Text = trimmed
			return compactContent(value)
		}
	}
	return compactContent(value)
}

func compactContent(value content.RichContent) content.RichContent {
	result := content.New()
	for _, block := range value.Blocks {
		nextBlock := content.Block{Type: block.Type}
		for _, child := range block.Children {
			if child.Type == "text" && child.Text == "" {
				continue
			}
			nextBlock.Children = append(nextBlock.Children, child)
		}
		if len(nextBlock.Children) > 0 {
			result.Blocks = append(result.Blocks, nextBlock)
		}
	}
	return result
}
