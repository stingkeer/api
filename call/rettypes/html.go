package rettypes

import (
	"html/template"
	"io"
	"io/fs"
	"reflect"
	"time"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/log"
)

var (
	_ def.ContentType = (*Html)(nil)
	_ def.RetAdapter  = (*Html)(nil)
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
			h.closeIO(write)
		}
	}()
	temp, err := template.New(h.tpl).Parse(h.tpl)
	if err != nil {
		log.Error(err)
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- temp.Execute(write, h.data)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error(err)
		}
	case <-time.After(time.Minute):
		log.Errorf("template execution timed out")
	}
	h.closeIO(write)
}

func (h *Html) closeIO(i any) {
	if w, b := i.(io.Closer); b {
		w.Close()
	}
}
