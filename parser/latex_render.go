package parser

import (
	"bytes"
	"github.com/russross/blackfriday"
)

var listing_lang map[string]bool

func init() {
	lang := []string{"abap", "acsl", "ada", "algol", "ant", "assembler",
		"awk", "bash", "basic", "c++", "c", "caml", "clean", "cobol",
		"comal", "csh", "delphi", "eiffel", "elan", "erlang", "euphoria",
		"fortran", "gcl", "gnuplot", "haskell", "html", "idl", "inform",
		"java", "jvmis", "ksh", "lisp", "logo", "lua", "make", "mathematica1",
		"matlab", "mercury", "metapost", "miranda", "mizar", "ml",
		"modelica3", "modula-", "mupad", "nastran", "oberon-", "ocl",
		"octave", "oz", "pascal", "perl", "php", "pl/i", "plasm", "pov",
		"prolog", "promela", "python", "r", "reduce", "rexx", "rsl",
		"ruby", "s", "sas", "scilab", "sh", "shelxl", "simula", "sql",
		"tcl", "tex", "vbscript", "verilog", "vhdl", "vrml", "xml", "xslt"}
	listing_lang = make(map[string]bool, len(lang))
	for _, c := range lang {
		listing_lang[c] = true
	}
}

// Latex is a type that implements the Renderer interface for LaTeX output.
//
// Do not create this directly, instead use the LatexRenderer function.
type Latex struct {
	PageLevel int
	Path      string
}

// LatexRenderer creates and configures a Latex object, which
// satisfies the Renderer interface.
//
// flags is a set of LATEX_* options ORed together (currently no such options
// are defined).
func LatexRenderer(flags int) blackfriday.Renderer {
	return &Latex{}
}

func (options *Latex) GetFlags() int {
	return 0
}

// render code chunks using verbatim, or listings if we have a language
func (options *Latex) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if _, ok := listing_lang[lang]; ok {
		out.WriteString("\n\\begin{lstlisting}[language=")
		out.WriteString(lang)
		out.WriteString("]\n")
		out.Write(text)
		out.WriteString("\n\\end{lstlisting}\n")
	} else {
		out.WriteString("\n\\begin{verbatim}\n")
		out.Write(text)
		out.WriteString("\n\\end{verbatim}\n")
	}
}

func (options *Latex) TitleBlock(out *bytes.Buffer, text []byte) {

}

func (options *Latex) BlockQuote(out *bytes.Buffer, text []byte) {
	out.WriteString("\n\\begin{quotation}\n")
	out.Write(text)
	out.WriteString("\n\\end{quotation}\n")
}

func (options *Latex) BlockHtml(out *bytes.Buffer, text []byte) {
	// a pretty lame thing to do...
	out.WriteString("\n\\begin{verbatim}\n")
	out.Write(text)
	out.WriteString("\n\\end{verbatim}\n")
}

func (options *Latex) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()

	switch level + options.PageLevel {
	case 1:
		out.WriteString("\n\\chapter{")
	case 2:
		out.WriteString("\n\\section{")
	case 3:
		out.WriteString("\n\\subsection{")
	case 4:
		out.WriteString("\n\\subsubsection{")
	case 5:
		out.WriteString("\n\\paragraph{")
	case 6:
		out.WriteString("\n\\subparagraph{")
	case 7:
		out.WriteString("\n\\textbf{")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("}\n")
}

func (options *Latex) HRule(out *bytes.Buffer) {
	out.WriteString("\n\\HRule\n")
}

func (options *Latex) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n\\begin{enumerate}\n")
	} else {
		out.WriteString("\n\\begin{itemize}\n")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n\\end{enumerate}\n")
	} else {
		out.WriteString("\n\\end{itemize}\n")
	}
}

func (options *Latex) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("\n\\item ")
	out.Write(text)
}

func (options *Latex) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	out.WriteString("\n")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}

func (options *Latex) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	out.WriteString("\n\\begin{tabular}{")
	for _, elt := range columnData {
		switch elt {
		case blackfriday.TABLE_ALIGNMENT_CENTER:
			out.WriteByte('c')
		case blackfriday.TABLE_ALIGNMENT_RIGHT:
			out.WriteByte('r')
		default:
			out.WriteByte('l')
		}
	}
	out.WriteString("}\n")
	out.Write(header)
	out.WriteString(" \\\\\n\\hline\n")
	out.Write(body)
	out.WriteString("\n\\end{tabular}\n")
}

func (options *Latex) TableRow(out *bytes.Buffer, text []byte) {
	if out.Len() > 0 {
		out.WriteString(" \\\\\n")
	}
	out.Write(text)
}

func (options *Latex) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	out.Write(text)
}

func (options *Latex) TableCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	escapeSpecialChars(out, text)
}

// TODO: this
func (options *Latex) Footnotes(out *bytes.Buffer, text func() bool) {

}

func (options *Latex) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {

}

func (options *Latex) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString("\\url{")
	if kind == blackfriday.LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	out.Write(link)
	out.WriteString("}")
}

func (options *Latex) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("\\texttt{")
	escapeSpecialChars(out, text)
	out.WriteString("}")
}

func (options *Latex) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\textbf{")
	escapeSpecialChars(out, text)
	out.WriteString("}")
}

func (options *Latex) Emphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\textit{")
	escapeSpecialChars(out, text)
	out.WriteString("}")
}

func (options *Latex) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	if bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://")) {
		// treat it like a link
		out.WriteString("\\href{")
		out.Write(bytes.Replace(link, []byte("}"), []byte("\\}"), -1))
		out.WriteString("}{")
		escapeSpecialChars(out, alt)
		out.WriteString("}")
	} else { //todo: pdf rewrite
		out.WriteString("\\includegraphics{")
		out.WriteString(options.Path)
		out.Write(bytes.Replace(link, []byte("}"), []byte("\\}"), -1))
		out.WriteString("}")
	}
}

func (options *Latex) LineBreak(out *bytes.Buffer) {
	out.WriteString(" \\\\\n")
}

func (options *Latex) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if len(link) > 0 && link[0] != '#' {
		out.WriteString("\\href{")
		out.Write(bytes.Replace(link, []byte("}"), []byte("\\}"), -1))
		out.WriteString("}{")
		escapeSpecialChars(out, content)
		out.WriteString("}")
	} else {
		escapeSpecialChars(out, content)
	}
}

func (options *Latex) RawHtmlTag(out *bytes.Buffer, tag []byte) {
}

func (options *Latex) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\textbf{\\textit{")
	escapeSpecialChars(out, text)
	out.WriteString("}}")
}

func (options *Latex) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.WriteString("\\sout{")
	escapeSpecialChars(out, text)
	out.WriteString("}")
}

// TODO: this
func (options *Latex) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {

}

func needsBackslash(c byte) bool {
	for _, r := range []byte("_{}%$&\\~#^") {
		if c == r {
			return true
		}
	}
	return false
}

func escapeSpecialChars(out *bytes.Buffer, text []byte) {
	for i := 0; i < len(text); i++ {
		// directly copy normal characters
		org := i

		for i < len(text) && !needsBackslash(text[i]) {
			i++
		}
		if i > org {
			out.Write(text[org:i])
		}

		// escape a character
		if i >= len(text) {
			break
		}
		out.WriteByte('\\')
		out.WriteByte(text[i])
	}
}

func (options *Latex) Entity(out *bytes.Buffer, entity []byte) {
	// TODO: convert this into a unicode character or something
	out.Write(entity)
}

func (options *Latex) NormalText(out *bytes.Buffer, text []byte) {
	escapeSpecialChars(out, text)
}

// header and footer
func (options *Latex) DocumentHeader(out *bytes.Buffer) {
}

func (options *Latex) DocumentFooter(out *bytes.Buffer) {
}

