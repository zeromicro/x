package http

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type message struct {
	XMLName xml.Name `json:"-" xml:"data"`
	Name    string   `json:"name" xml:"name"`
}

func TestOkXml(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	OkXml(&w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "<data><name>anyone</name></data>", w.builder.String())
}

func TestOkXmlCtx(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	OkXmlCtx(context.TODO(), &w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "<data><name>anyone</name></data>", w.builder.String())
}

func TestWriteXmlTimeout(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     http.ErrHandlerTimeout,
	}
	msg := message{Name: "anyone"}
	WriteXml(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteXmlError(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     errors.New("foo"),
	}
	msg := message{Name: "anyone"}
	WriteXml(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteXmlCtxError(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     errors.New("foo"),
	}
	msg := message{Name: "anyone"}
	WriteXmlCtx(context.TODO(), &w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteXmlLessWritten(t *testing.T) {
	w := tracedResponseWriter{
		headers:     make(map[string][]string),
		lessWritten: true,
	}
	msg := message{Name: "anyone"}
	WriteXml(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteXmlMarshalFailed(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	WriteXml(&w, http.StatusOK, map[string]any{
		"Data": complex(0, 0),
	})
	assert.Equal(t, http.StatusInternalServerError, w.code)
}

type tracedResponseWriter struct {
	headers     map[string][]string
	builder     strings.Builder
	hasBody     bool
	code        int
	lessWritten bool
	wroteHeader bool
	err         error
}

type testWriterResult struct {
	code        int
	writeString string
}

func (w *tracedResponseWriter) result() (testWriterResult, error) {
	return testWriterResult{
		code:        w.code,
		writeString: w.builder.String(),
	}, w.err
}

func (w *tracedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *tracedResponseWriter) Write(bytes []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err = w.builder.Write(bytes)
	if w.lessWritten {
		n--
	}
	w.hasBody = true

	return
}

func (w *tracedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.code = code
}

func TestOkHTML(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := `<!DOCTYPE html>
<html>
<head>
  <title>Hello, World!</title>
</head>
<body>
  <h1>Hello, World!</h1>
</body>
</html>`
	OkHTML(&w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "<!DOCTYPE html>\n<html>\n<head>\n  <title>Hello, World!</title>\n</head>\n<body>\n  <h1>Hello, World!</h1>\n</body>\n</html>", w.builder.String())
}
func TestOkHTMLCtx(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := `<!DOCTYPE html>
<html>
<head>
  <title>Hello, World!</title>
</head>
<body>
  <h1>Hello, World!</h1>
</body>
</html>`
	OkHTMLCtx(context.TODO(), &w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "<!DOCTYPE html>\n<html>\n<head>\n  <title>Hello, World!</title>\n</head>\n<body>\n  <h1>Hello, World!</h1>\n</body>\n</html>", w.builder.String())
}

func TestWriteHTMLTimeout(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     http.ErrHandlerTimeout,
	}
	msg := `<!DOCTYPE html>
<html>
<head>
  <title>Hello, World!</title>
</head>
<body>
  <h1>Hello, World!</h1>
</body>
</html>`
	WriteHTML(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteHTMLError(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     errors.New("foo"),
	}
	msg := `<!DOCTYPE html>
<html>
<head>
  <title>Hello, World!</title>
</head>
<body>
  <h1>Hello, World!</h1>
</body>
</html>`
	WriteHTML(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteHTMLCtxError(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     errors.New("foo"),
	}
	msg := `<!DOCTYPE html>
<html>
<head>
  <title>Hello, World!</title>
</head>
<body>
  <h1>Hello, World!</h1>
</body>
</html>`
	WriteHTMLCtx(context.TODO(), &w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteHTMLLessWritten(t *testing.T) {
	w := tracedResponseWriter{
		headers:     make(map[string][]string),
		lessWritten: true,
	}
	msg := `<!DOCTYPE html>
<html>
<head>
  <title>Hello, World!</title>
</head>
<body>
  <h1>Hello, World!</h1>
</body>
</html>`
	WriteHTML(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

type MyString string

func TestWritHTMLTypeFailed(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	WriteHTML(&w, http.StatusOK, map[string]any{
		"Data": complex(0, 0),
	})
	assert.Equal(t, http.StatusInternalServerError, w.code)

	w = tracedResponseWriter{
		headers: make(map[string][]string),
	}
	WriteHTML(&w, http.StatusOK, MyString("foo"))
	assert.Equal(t, http.StatusInternalServerError, w.code)
}
