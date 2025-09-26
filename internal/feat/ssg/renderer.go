package ssg

import (
	"fmt"

	gmast "github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// TailwindRenderer is a custom renderer for goldmark that adds Tailwind CSS classes.
type TailwindRenderer struct {
	html.Config
}

// NewTailwindRenderer creates a new TailwindRenderer.
func NewTailwindRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TailwindRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers the render functions for the nodes.
func (r *TailwindRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gmast.KindHeading, r.renderHeading)
	reg.Register(gmast.KindParagraph, r.renderParagraph)
	reg.Register(gmast.KindList, r.renderList)
	reg.Register(gmast.KindListItem, r.renderListItem)
	reg.Register(gmast.KindBlockquote, r.renderBlockquote)
	reg.Register(gmast.KindThematicBreak, r.renderHorizontalRule)
	reg.Register(gmast.KindEmphasis, r.renderEmphasis)
	reg.Register(extast.KindStrikethrough, r.renderDel)
	reg.Register(gmast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(gmast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(gmast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(extast.KindTable, r.renderTable)
	reg.Register(extast.KindTableHeader, r.renderTableHeader)
	reg.Register(extast.KindTableRow, r.renderTableRow)
	reg.Register(extast.KindTableCell, r.renderTableCell)
	reg.Register(gmast.KindLink, r.renderLink)
	reg.Register(gmast.KindImage, r.renderImage)
}

func (r *TailwindRenderer) renderHeading(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Heading)
	if entering {
		level := n.Level
		var class string
		switch level {
		case 1:
			class = "prose-h1"
		case 2:
			class = "prose-h2"
		case 3:
			class = "prose-h3"
		case 4:
			class = "prose-h4"
		case 5:
			class = "prose-h5"
		case 6:
			class = "prose-h6"
		}
		_, _ = w.WriteString(fmt.Sprintf("<h%d class=\"%s\">", level, class))
	} else {
		_, _ = w.WriteString(fmt.Sprintf("</h%d>", n.Level))
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderParagraph(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<p class=\"prose-p\">")
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderList(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.List)
	tag := "ul"
	if n.IsOrdered() {
		tag = "ol"
	}
	if entering {
		_, _ = w.WriteString(fmt.Sprintf("<%s class=\"prose-ul\">", tag))
	} else {
		_, _ = w.WriteString(fmt.Sprintf("</%s>", tag))
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderListItem(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<li class=\"prose-li\">")
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderBlockquote(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<blockquote class=\"prose-blockquote\">")
	} else {
		_, _ = w.WriteString("</blockquote>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderHorizontalRule(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<hr class=\"prose-hr\">\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderEmphasis(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Emphasis)
	if entering {
		if n.Level == 2 {
			_, _ = w.WriteString("<strong>")
		} else {
			_, _ = w.WriteString("<em>")
		}
	} else {
		if n.Level == 2 {
			_, _ = w.WriteString("</strong>")
		} else {
			_, _ = w.WriteString("</em>")
		}
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderDel(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<del>")
	} else {
		_, _ = w.WriteString("</del>")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderCodeSpan(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<code class=\"prose-code\">")
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*gmast.Text).Segment
			_, _ = w.Write(segment.Value(source))
		}
		_, _ = w.WriteString("</code>")
	}
	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderCodeBlock(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre class=\"prose-pre\"><code>")
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.Write(line.Value(source))
		}
		_, _ = w.WriteString("</code></pre>\n")
	}
	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre class=\"prose-pre\"><code>")
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.Write(line.Value(source))
		}
		_, _ = w.WriteString("</code></pre>\n")
	}
	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderTable(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<table class=\"prose-table\">")
	} else {
		_, _ = w.WriteString("</table>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderTableHeader(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<thead>")
	} else {
		_, _ = w.WriteString("</thead>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderTableRow(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<tr>")
	} else {
		_, _ = w.WriteString("</tr>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderTableCell(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<td>")
	} else {
		_, _ = w.WriteString("</td>")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderLink(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Link)
	if entering {
		_, _ = w.Write([]byte(fmt.Sprintf("<a href=\"%s\" class=\"prose-a\">", n.Destination)))
	} else {
		_, _ = w.Write([]byte("</a>"))
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderImage(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Image)
	if entering {
		_, _ = w.Write([]byte(fmt.Sprintf("<img src=\"%s\" alt=\"", n.Destination)))
		// The alt text is a child of the image node.
	} else {
		_, _ = w.Write([]byte("\" class=\"prose-img\">"))
	}
	return gmast.WalkContinue, nil
}
