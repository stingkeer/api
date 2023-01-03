package rettypes

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"
	"html/template"
	"io"
	"io/fs"
	"reflect"
)

type Html struct {
	tpl    string
	data   any
	reader io.Reader
	writer io.Writer
}

func NewHtml(tpl string, data any) *Html {
	reader, writer := io.Pipe()
	h := &Html{tpl: tpl, data: data, reader: reader, writer: writer}
	return h
}

func HtmlView(root fs.FS, name string, data any) *Html {
	reader, writer := io.Pipe()
	open, err := root.Open(name)
	if err != nil {
		panic(err)
	}
	all, err := io.ReadAll(open)
	if err != nil {
		panic(err)
	}
	h := &Html{data: data, reader: reader, writer: writer, tpl: string(all)}
	return h
}

func (h *Html) ContentType() string {
	return def.CONTENT_HTML
}

func (h *Html) Return() io.Reader {
	go h.renderer(h.writer)
	return h.reader
}

func (h *Html) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Html)(nil))}
}

func (h *Html) renderer(write io.Writer) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			if w, b := write.(*io.PipeWriter); b {
				w.Close()
			}
		}
	}()
	temp, err := template.New("aaa").Parse(h.tpl)
	if err != nil {
		panic(err)
	}
	err = temp.Execute(write, h.data)
	if err != nil {
		panic(err)
	}
	if w, b := write.(*io.PipeWriter); b {
		w.Close()
	}
}
