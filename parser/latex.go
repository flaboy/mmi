package parser

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func (n *Node) ToLatex() {
	n.to_latex(os.Stdout, &Latex{}, 0)
}

func (n *Node) to_latex(out io.Writer, r *Latex, depth int) {

	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	depth += 1
	for _, c := range n.Child {
		r.PageLevel = depth
		if c.IsPage == false {
			buf := bytes.NewBufferString("")
			r.Header(buf, func() bool { buf.WriteString(c.Title); return true }, 1, "")
			out.Write(buf.Bytes())
			c.to_latex(out, r, depth)
		} else {
			fd, err := os.OpenFile(c.filepath, os.O_RDONLY, 0644)
			if err == nil {
				buf, err := ioutil.ReadAll(fd)
				fd.Close()

				if err != nil {
					panic(err)
				}
				fmt.Fprintln(os.Stderr, n.workdir)
				r.Path = n.workdir + "/" + strings.Join(n.Path, "/")
				if r.Path != "" {
					r.Path += "/"
				}

				if len(buf) > 0 {
					if buf[0] != '#' {
						buf = append([]byte("# "), buf...)
					}

					fmt.Fprintf(out, "%s\n", blackfriday.Markdown(buf, r, extensions))
				}

			} else {
				fmt.Println(err)
			}
		}
	}
}
