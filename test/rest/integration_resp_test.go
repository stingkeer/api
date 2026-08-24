package rest

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
)

// TestIntegrationRespCustomCodeAndHeaders: NewResp with status code and
// custom headers.
func TestIntegrationRespCustomCodeAndHeaders(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.POST(func() any {
			return api.NewResp(map[string]string{"created": "yes"}).
				SetCode(http.StatusCreated).
				SetHeader(map[string]string{"X-Trace-Id": "abc-123"})
		}, "/resp/created")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusCreated)
		resp.AssertHeader("X-Trace-Id", "abc-123")
		resp.AssertHeader("Content-Type", def.Content_JSON)
		resp.AssertBody(`{"created":"yes"}`)
	})
}

// TestIntegrationRespRawReader: NewResp.SetReader streams the reader as-is
// with a custom content type.
func TestIntegrationRespRawReader(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewResp(nil).
				SetReader(strings.NewReader("raw-payload")).
				SetContentType("application/octet-stream")
		}, "/resp/raw")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertHeader("Content-Type", "application/octet-stream")
		resp.AssertBody("raw-payload")
	})
}

// TestIntegrationRedirect: NewRedirect answers 302 with a Location header.
func TestIntegrationRedirect(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewRedirect("https://example.com/target")
		}, "/resp/redirect")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusFound)
		resp.AssertHeader("Location", "https://example.com/target")
	})
}

// TestIntegrationStreamDownload: NewStream serves bytes from memory; with
// SetName the response carries a Content-Disposition attachment header.
func TestIntegrationStreamDownload(t *testing.T) {
	const payload = "0123456789abcdef-STREAM-CONTENT"
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewStream(strings.NewReader(payload)).SetName("data.bin")
		}, "/resp/download")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		if cd := resp.Header("Content-Disposition"); !strings.Contains(cd, "data.bin") {
			t.Errorf("Content-Disposition: expect attachment name data.bin but %q", cd)
		}
		resp.AssertBody(payload)
	})
}

// TestIntegrationPanicHandler: a panicking handler becomes a JSON 500 error,
// and the server keeps serving afterwards.
func TestIntegrationPanicHandler(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { panic("integration-boom") }, "/resp/panic")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusInternalServerError)
		var e def.Error
		if err := jsonUnmarshalString(resp.BodyString(), &e); err != nil {
			t.Fatal(err)
		}
		if e.ErrorMessage != "integration-boom" {
			t.Errorf("error message: expect integration-boom but %q", e.ErrorMessage)
		}
	})
	// server survived the panic
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "alive" }, "/resp/after-panic")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody("alive")
	})
}

// TestIntegrationNilAndEmptyReturns: nil and no-return handlers answer 200
// with an empty body.
func TestIntegrationNilAndEmptyReturns(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { return nil }, "/resp/nil")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody("")
	})
	r.Test(t, func() def.Option {
		return api.GET(func() {}, "/resp/empty")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody("")
	})
}

// TestIntegrationSerializedCollections: maps, slices and structs serialize
// to JSON with the JSON content type.
func TestIntegrationSerializedCollections(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { return []int{3, 1, 2} }, "/resp/slice")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertHeader("Content-Type", def.Content_JSON)
		resp.AssertBody("[3,1,2]")
	})
	type item struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	r.Test(t, func() def.Option {
		return api.GET(func() any { return item{A: "x", B: 9} }, "/resp/struct")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertBody(`{"a":"x","b":9}`)
	})
}

// TestIntegrationConcurrentRequests: many parallel clients against one
// handler — every response must be correct and no request may be lost.
func TestIntegrationConcurrentRequests(t *testing.T) {
	var hits atomic.Int32
	r.Test(t, func() def.Option {
		return api.GET(func(n string) any {
			hits.Add(1)
			return "n=" + n
		}, "/resp/concurrent")
	})

	const workers, perWorker = 16, 12
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				resp, err := http.Get(r.BaseURL() + "/resp/concurrent?n=" + itoa(w*perWorker+i))
				if err != nil {
					errCh <- err
					return
				}
				body, _ := ioReadAll(resp)
				resp.Body.Close()
				if want := "n=" + itoa(w*perWorker+i); string(body) != want {
					errCh <- &mismatchError{want: want, got: string(body)}
					return
				}
			}
		}(w)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}
	if got := hits.Load(); got != workers*perWorker {
		t.Errorf("handler hits: expect %d but %d", workers*perWorker, got)
	}
}

type mismatchError struct{ want, got string }

func (e *mismatchError) Error() string {
	return "concurrent body: expect " + e.want + " but " + e.got
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
