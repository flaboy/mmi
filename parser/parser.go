package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
)

type (
	pathway []string
	Node    struct {
		filepath string
		Title    string
		Path     pathway
		PathNode []*Node
		IsPage   bool
		Child    []*Node
		next     *Node
		prev     *Node
	}
)

var current_node *Node

const (
	readme_md  = "README.md"
	summary_md = "SUMMARY.md"
)

func Open(dirname string) Node {
	return opendir(dirname, pathway{}, []*Node{})
}

func (n *Node) UpdateRummary(depth int64) {
	n.update_summary(depth, 0)
}

func (n *Node) link(target *Node) interface{} {
	if target == nil {
		return nil
	} else {
		link := strings.Join(target.Path, "/")
		if target.IsPage == false {
			link += "/README.md"
		}
		return map[string]string{
			"path":  link,
			"title": target.Title,
		}
	}
}

func (n *Node) UpdateJson() {
	fmt.Println("[U]", n.filepath+"/index.json")
	f, err := os.OpenFile(n.filepath+"/index.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		panic(err)
	}

	child := map[string]interface{}{
		"README.md": map[string]interface{}{
			"title": n.Title,
			"next":  n.link(n.next),
			"prev":  n.link(n.prev),
		},
	}
	for _, c := range n.Child {
		if c.IsPage {
			child[path.Base(c.filepath)] = map[string]interface{}{
				"title": c.Title,
				"next":  c.link(c.next),
				"prev":  c.link(c.prev),
			}
		}
	}
	data := map[string]interface{}{
		"nav":   n.navpath(),
		"child": child,
	}
	b, err := json.MarshalIndent(data, "", "    ")
	f.Write(b)
	f.Close()

	for _, c := range n.Child {
		if c.IsPage == false {
			c.UpdateJson()
		}
	}
}

func (c *Node) navpath() []interface{} {
	links := []interface{}{}
	for _, n := range c.PathNode {
		links = append(links, c.link(n))
	}
	links = append(links, c.link(c))
	return links
}

func (n *Node) update_summary(depth, st int64) {
	fmt.Println("[U]", n.filepath+"/"+summary_md)
	f, err := os.OpenFile(n.filepath+"/"+summary_md, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}

	f.Truncate(0)
	f.WriteString(n.Title)
	f.WriteString("\n================================================\n\n")
	f.Write([]byte(n.TocMarkdown(depth, st, "")))
	f.Close()

	for _, c := range n.Child {
		if c.IsPage == false {
			c.update_summary(depth, st+1)
		}
	}
}

func (n *Node) TocMarkdown(depth, st int64, prefix string) (s string) {
	tab := "    "
	for _, c := range n.Child {
		filepath := strings.Join(c.Path[st:], "/")
		if c.IsPage == false {
			filepath += "/" + readme_md
		}

		s += fmt.Sprintf("%s1. [%s](%s)\n", prefix, c.Title, filepath)

		if c.IsPage == false {
			s += c.TocMarkdown(depth-1, st, prefix+tab)
		}
	}
	return
}

func link_node(n *Node) {
	if current_node == nil {
		current_node = n
	} else {
		current_node.next = n
		n.prev = current_node
		current_node = n
	}
}

func opendir(dirname string, p pathway, pn []*Node) Node {
	flist, err := ioutil.ReadDir(dirname)

	n := Node{Path: p, filepath: dirname}
	link_node(&n)
	pn = append(pn, &n)

	n.Title = get_title(dirname + "/" + readme_md)

	var subnode *Node
	if err == nil {
		for _, f := range flist {
			fname := f.Name()
			fpath := dirname + "/" + fname
			subpath := append(p, url.QueryEscape(fname))
			subnode = nil
			if fname[0] != '.' && fname != readme_md && fname != summary_md {
				if f.Mode().IsDir() {
					_, err = os.Stat(fpath + "/" + readme_md)
					if err == nil {
						new_node := opendir(fpath, subpath, pn)
						subnode = &new_node
					}
				} else if path.Ext(fname) == ".md" {
					new_node := openfile(fpath, subpath)
					subnode = &new_node
					link_node(subnode)
				}
			}
			if subnode != nil {
				subnode.PathNode = pn
				n.Child = append(n.Child, subnode)
			}
		}
	}
	return n
}

func openfile(file string, p pathway) Node {
	n := Node{Path: p, IsPage: true, filepath: file}
	n.Title = get_title(file)
	return n
}

func get_title(file string) (title string) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err == nil {
		r := bufio.NewReader(f)
		l, _ := r.ReadString('\n')
		if len(l) > 2 && l[0:2] == "# " {
			title = l[2:]
		} else {
			title = l
		}
	}
	return strings.Trim(title, "\n")
}
