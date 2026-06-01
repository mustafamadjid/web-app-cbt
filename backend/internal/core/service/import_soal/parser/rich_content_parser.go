package parser

import (
	"encoding/xml"
	"fmt"
	"strings"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
)

type xmlNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Nodes   []xmlNode  `xml:",any"`
	Text    string     `xml:",chardata"`
}

func ExtractParagraphContents(data []byte) ([]content.RichContent, []string, error) {
	documentBytes, err := readDocumentXML(data)
	if err != nil {
		return nil, nil, err
	}

	var root xmlNode
	if err := xml.Unmarshal(documentBytes, &root); err != nil {
		return nil, nil, fmt.Errorf("unmarshal document.xml: %w", err)
	}

	bodyNode := findFirstNode(root, "body")
	if bodyNode == nil {
		return nil, nil, fmt.Errorf("document body not found in docx")
	}

	paragraphs := make([]content.RichContent, 0)
	warnings := make([]string, 0)
	for _, node := range bodyNode.Nodes {
		if node.XMLName.Local != "p" {
			continue
		}

		rich, paragraphWarnings := paragraphNodeToContent(node)
		if !rich.Empty() {
			paragraphs = append(paragraphs, rich)
		}
		warnings = append(warnings, paragraphWarnings...)
	}

	return paragraphs, warnings, nil
}

func paragraphNodeToContent(node xmlNode) (content.RichContent, []string) {
	result := content.New()
	block := content.Block{Type: "paragraph"}
	warnings := make([]string, 0)

	appendText := func(text string, marks []content.Mark) {
		if text == "" {
			return
		}
		text = strings.ReplaceAll(text, "\u00a0", " ")
		block.Children = append(block.Children, content.Inline{
			Type:  "text",
			Text:  text,
			Marks: cloneMarks(marks),
		})
	}

	var visitNode func(current xmlNode, inheritedMarks []content.Mark)
	visitNode = func(current xmlNode, inheritedMarks []content.Mark) {
		switch current.XMLName.Local {
		case "r":
			runMarks := mergeMarks(inheritedMarks, marksFromRun(current))
			for _, child := range current.Nodes {
				visitNode(child, runMarks)
			}
		case "t":
			appendText(current.Text, inheritedMarks)
		case "tab":
			appendText("\t", inheritedMarks)
		case "br":
			appendText("\n", inheritedMarks)
		case "oMath", "oMathPara":
			appendMathInline(&block, current, &warnings)
		default:
			for _, child := range current.Nodes {
				visitNode(child, inheritedMarks)
			}
		}
	}

	for _, child := range node.Nodes {
		visitNode(child, nil)
	}

	if len(block.Children) > 0 {
		result.Blocks = append(result.Blocks, block)
	}
	return result, warnings
}

func appendMathInline(block *content.Block, node xmlNode, warnings *[]string) {
	latex, warning := ommlNodeToLatex(node)
	if warning != "" {
		*warnings = append(*warnings, warning)
	}
	if strings.TrimSpace(latex) == "" {
		return
	}

	display := "inline"
	if node.XMLName.Local == "oMathPara" {
		display = "block"
	}
	block.Children = append(block.Children, content.Inline{
		Type:    "math",
		Latex:   latex,
		Display: display,
	})
}
