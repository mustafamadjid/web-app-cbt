package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func readDocumentXML(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open document.xml: %w", err)
		}
		defer rc.Close()
		xmlBytes, err := io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("read document.xml: %w", err)
		}
		return xmlBytes, nil
	}

	return nil, fmt.Errorf("word/document.xml not found in docx")
}

func findFirstNode(node xmlNode, local string) *xmlNode {
	if node.XMLName.Local == local {
		return &node
	}
	for _, child := range node.Nodes {
		if child.XMLName.Local == local {
			return &child
		}
		if nested := findFirstNode(child, local); nested != nil {
			return nested
		}
	}
	return nil
}

func currentTextValue(node xmlNode) string {
	var sb strings.Builder
	var visit func(xmlNode)
	visit = func(item xmlNode) {
		if item.XMLName.Local == "t" && item.Text != "" {
			sb.WriteString(item.Text)
		}
		if item.Text != "" && item.XMLName.Local == "" {
			sb.WriteString(item.Text)
		}
		for _, child := range item.Nodes {
			visit(child)
		}
	}
	visit(node)
	return sb.String()
}

func attrValue(node xmlNode, local string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Local == local {
			return attr.Value
		}
	}
	return ""
}

