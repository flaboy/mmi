package main

import (
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

func err_f(code int, w http.ResponseWriter, req *http.Request) {
}

func handler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)

	try_files := []string{"README.md",
		"index.md", "INDEX.md", "index.html", "index.htm"}

	filepath := workdir + req.URL.Path
	st, err := os.Stat(filepath)

	if err != nil {
		err_f(403, w, req)
		return
	}

	if st.IsDir() && req.URL.Path[len(req.URL.Path)-1:len(req.URL.Path)] != "/" {
		w.Header().Set("Location", req.URL.Path+"/")
		return
	}

	for _, tf := range try_files {
		test_file := filepath + "/" + tf
		_, err := os.Stat(test_file)
		if err == nil {
			filepath = test_file
			break
		}
	}

	fd, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	defer fd.Close()

	if err != nil {
		err_f(403, w, req)
		return
	}

	ext_name := path.Ext(filepath)
	if ext_name == ".md" {
		buf, err := ioutil.ReadAll(fd)
		if err != nil {
			err_f(403, w, req)
			return
		}
		html_code := to_html(buf)

		html_render(w, html_code)
		return
	}

	content_type := mime.TypeByExtension(ext_name)
	if content_type != "" {
		w.Header().Set("Content-Type", content_type)
	}
	io.Copy(w, fd)
}

func start_server(workdir string) {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":10080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func to_html(md_code []byte) []byte {
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	renderer := blackfriday.HtmlRenderer(htmlFlags, "title", "css")

	return blackfriday.Markdown(md_code, renderer, extensions)
}
