package rest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	ihttp "go.aew.app/api.v1/http"
	"go.aew.app/api.v1/test/r"
)

// TestIntegrationUploadMultipart: multipart/form-data bodies bind to a
// multipart.Reader parameter.
func TestIntegrationUploadMultipart(t *testing.T) {
	cl := r.Test(t, func() def.Option {
		return api.POST(func(read multipart.Reader) any {
			part, err := read.NextPart()
			if err != nil {
				return api.NewResp(map[string]string{"error": err.Error()}).SetCode(http.StatusBadRequest)
			}
			defer part.Close()
			b, _ := io.ReadAll(part)
			return map[string]any{
				"filename": part.FileName(),
				"size":     len(b),
				"content":  string(b),
			}
		}, "/files/upload")
	})

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "note.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write([]byte("upload payload")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	cl.Request().
		AddHeader("Content-Type", w.FormDataContentType()).
		SetBody(buf.Bytes()).
		Do(func(resp *r.Response) {
			resp.AssertStatus(http.StatusOK)
			resp.AssertBody(`{"content":"upload payload","filename":"note.txt","size":14}`)
		})
}

// TestIntegrationDownloadRange: NewStream over a real file serves Range
// requests as 206 with Content-Range and the requested byte slice.
// (The compressor deliberately skips ranged responses.)
func TestIntegrationDownloadRange(t *testing.T) {
	const payload = "0123456789abcdefghijklmnopqrstuvwxyz"
	file := filepath.Join(t.TempDir(), "range.bin")
	if err := os.WriteFile(file, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}

	r.Test(t, func() def.Option {
		return api.GET(func() any {
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			return api.NewStream(f).SetName("range.bin")
		}, "/files/range")
	}).Request().AddHeader("Range", "bytes=0-9").Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusPartialContent)
		resp.AssertHeader("Content-Range", "bytes 0-9/"+itoa(len(payload)))
		resp.AssertHeader("Accept-Ranges", "bytes")
		resp.AssertBody(payload[:10])
	})
}

// TestIntegrationStaticFiles: api.Static serves files from an http.FileSystem;
// unknown files fall through to 404.
func TestIntegrationStaticFiles(t *testing.T) {
	dir := t.TempDir()
	const content = "static hello content"
	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	api.Static("/web/*", "", http.Dir(dir), ihttp.StaticRewrite("/web/", ""))

	// boot the shared server (r.Test) before issuing raw requests
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "boot" }, "/files/boot")
	})

	resp, err := http.Get(r.BaseURL() + "/web/hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ioReadAll(resp)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(b) != content {
		t.Errorf("static file: code=%d body=%q", resp.StatusCode, string(b))
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("static Content-Type: expect text/plain but %q", ct)
	}

	// missing file → 404 (falls through the static handler)
	resp404, err := http.Get(r.BaseURL() + "/web/missing.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer resp404.Body.Close()
	if resp404.StatusCode != http.StatusNotFound {
		t.Errorf("missing static file: expect 404 but %d", resp404.StatusCode)
	}
}
