package parser

import (
	"fmt"
	"strings"
)

func ommlNodeToLatex(node xmlNode) (string, string) {
	switch node.XMLName.Local {
	case "oMathPara":
		return ommlJoinedLatex(node.Nodes, " ")
	case "oMath":
		return ommlJoinedLatex(node.Nodes, "")
	case "r":
		return gatherOMMLText(node), ""
	case "t":
		return currentTextValue(node), ""
	case "sSup":
		return ommlScriptLatex(node, "superscript", func(base, sup string) string {
			return fmt.Sprintf("{%s}^{%s}", base, sup)
		}, "sup")
	case "sSub":
		return ommlScriptLatex(node, "subscript", func(base, sub string) string {
			return fmt.Sprintf("{%s}_{%s}", base, sub)
		}, "sub")
	case "sSubSup":
		return ommlSubSupLatex(node)
	case "f":
		return ommlFractionLatex(node)
	case "rad":
		return ommlRadicalLatex(node)
	case "d":
		return ommlDelimiterLatex(node)
	case "nary":
		return ommlNaryLatex(node)
	case "acc":
		return ommlAccentLatex(node)
	case "m":
		return ommlMatrixLatex(node)
	case "func":
		name := ommlChildLatex(node, "fName")
		arg := ommlChildLatex(node, "e")
		return name + "{" + arg + "}", ""
	default:
		text := gatherOMMLText(node)
		if strings.TrimSpace(text) == "" {
			return "", ""
		}
		return text, fmt.Sprintf("unsupported OMML node %q fallback ke raw text", node.XMLName.Local)
	}
}

func ommlJoinedLatex(nodes []xmlNode, separator string) (string, string) {
	var parts []string
	var warnings []string
	for _, child := range nodes {
		latex, warning := ommlNodeToLatex(child)
		if latex != "" {
			parts = append(parts, latex)
		}
		if warning != "" {
			warnings = append(warnings, warning)
		}
	}
	return strings.TrimSpace(strings.Join(parts, separator)), strings.Join(warnings, "; ")
}

func ommlScriptLatex(node xmlNode, label string, format func(string, string) string, childName string) (string, string) {
	base := ommlChildLatex(node, "e")
	child := ommlChildLatex(node, childName)
	if base == "" || child == "" {
		return gatherOMMLText(node), fmt.Sprintf("unsupported OMML %s fallback ke raw text", label)
	}
	return format(base, child), ""
}

func ommlSubSupLatex(node xmlNode) (string, string) {
	base := ommlChildLatex(node, "e")
	sub := ommlChildLatex(node, "sub")
	sup := ommlChildLatex(node, "sup")
	if base == "" {
		return gatherOMMLText(node), "unsupported OMML subscript/superscript fallback ke raw text"
	}
	return fmt.Sprintf("{%s}_{%s}^{%s}", base, sub, sup), ""
}

func ommlFractionLatex(node xmlNode) (string, string) {
	num := ommlChildLatex(node, "num")
	den := ommlChildLatex(node, "den")
	if num == "" || den == "" {
		return gatherOMMLText(node), "unsupported OMML fraction fallback ke raw text"
	}
	return fmt.Sprintf("\\frac{%s}{%s}", num, den), ""
}

func ommlRadicalLatex(node xmlNode) (string, string) {
	deg := ommlChildLatex(node, "deg")
	base := ommlChildLatex(node, "e")
	if base == "" {
		return gatherOMMLText(node), "unsupported OMML radical fallback ke raw text"
	}
	if deg != "" {
		return fmt.Sprintf("\\sqrt[%s]{%s}", deg, base), ""
	}
	return fmt.Sprintf("\\sqrt{%s}", base), ""
}

func ommlDelimiterLatex(node xmlNode) (string, string) {
	open, close := "(", ")"
	if beg := findFirstNode(node, "begChr"); beg != nil {
		if v := attrValue(*beg, "val"); v != "" {
			open = v
		}
	}
	if end := findFirstNode(node, "endChr"); end != nil {
		if v := attrValue(*end, "val"); v != "" {
			close = v
		}
	}
	body := ommlChildLatex(node, "e")
	return fmt.Sprintf("\\left%s %s \\right%s", open, body, close), ""
}

func ommlNaryLatex(node xmlNode) (string, string) {
	body := ommlChildLatex(node, "e")
	sub := ommlChildLatex(node, "sub")
	sup := ommlChildLatex(node, "sup")
	op := "\\sum"
	if chr := findFirstNode(node, "chr"); chr != nil {
		switch attrValue(*chr, "val") {
		case "\u220f":
			op = "\\prod"
		case "\u22c2":
			op = "\\bigcap"
		case "\u22c3":
			op = "\\bigcup"
		}
	}
	return op + optionalBound("_", sub) + optionalBound("^", sup) + "{" + body + "}", ""
}

func ommlAccentLatex(node xmlNode) (string, string) {
	base := ommlChildLatex(node, "e")
	if base == "" {
		return gatherOMMLText(node), "unsupported OMML accent fallback ke raw text"
	}

	cmd := "\\hat"
	if chr := findFirstNode(node, "chr"); chr != nil {
		switch attrValue(*chr, "val") {
		case "\u00af":
			cmd = "\\bar"
		case "\u2192":
			cmd = "\\vec"
		case "\u02d9":
			cmd = "\\dot"
		}
	}
	return fmt.Sprintf("%s{%s}", cmd, base), ""
}

func ommlMatrixLatex(node xmlNode) (string, string) {
	rows := make([]string, 0)
	for _, row := range node.Nodes {
		if row.XMLName.Local != "mr" {
			continue
		}
		cells := make([]string, 0)
		for _, cell := range row.Nodes {
			if cell.XMLName.Local == "e" {
				cells = append(cells, ommlChildrenLatex(cell.Nodes))
			}
		}
		rows = append(rows, strings.Join(cells, " & "))
	}
	if len(rows) == 0 {
		return gatherOMMLText(node), "unsupported OMML matrix fallback ke raw text"
	}
	return "\\begin{matrix}" + strings.Join(rows, ` \\ `) + "\\end{matrix}", ""
}

func gatherOMMLText(node xmlNode) string {
	var sb strings.Builder
	var visit func(xmlNode)
	visit = func(item xmlNode) {
		if item.Text != "" {
			sb.WriteString(item.Text)
		}
		for _, child := range item.Nodes {
			visit(child)
		}
	}
	visit(node)
	return strings.TrimSpace(sb.String())
}

func ommlChildLatex(node xmlNode, local string) string {
	for _, child := range node.Nodes {
		if child.XMLName.Local == local {
			return ommlChildrenLatex(child.Nodes)
		}
	}
	return ""
}

func ommlChildrenLatex(nodes []xmlNode) string {
	parts := make([]string, 0, len(nodes))
	for _, child := range nodes {
		latex, _ := ommlNodeToLatex(child)
		if latex != "" {
			parts = append(parts, latex)
		}
	}
	return strings.Join(parts, "")
}

func optionalBound(prefix, value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return prefix + "{" + value + "}"
}
