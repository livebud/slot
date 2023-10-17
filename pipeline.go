package slot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func newPipeline(rw http.ResponseWriter) *pipeline {
	r, wc := newPipe()
	return &pipeline{
		rw: rw,
		state: &state{
			Named: map[string]string{},
		},
		headers: http.Header{},
		reader:  bytes.NewBuffer(nil),
		writer:  wc,
		next:    r,
	}
}

type pipeline struct {
	rw      http.ResponseWriter
	state   *state
	headers http.Header
	reader  io.Reader
	writer  io.WriteCloser
	next    io.Reader
}

var _ http.ResponseWriter = (*pipeline)(nil)

func (pl *pipeline) Header() http.Header {
	return pl.headers
}

func (pl *pipeline) WriteHeader(statusCode int) {
	pl.rw.WriteHeader(statusCode)
}

func (pl *pipeline) Write(p []byte) (n int, err error) {
	pl.state.Data += string(p)
	return len(p), nil
}

func (pl *pipeline) WriteTo(slot string, p []byte) (n int, err error) {
	pl.state.Named[slot] += string(p)
	return len(p), nil
}

func (pl *pipeline) Close() error {
	enc := json.NewEncoder(pl.writer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(pl.state); err != nil {
		return err
	}
	// Flush headers as we close the writer
	rh := pl.rw.Header()
	for name := range pl.headers {
		value := pl.headers.Get(name)
		// Prefer headers set earlier in the pipeline vs later, they're more
		// specific, e.g. view vs. layout.
		if rh.Get(name) == "" {
			rh.Set(name, value)
		}
	}
	return pl.writer.Close()
}

func (pl *pipeline) Read(p []byte) (n int, err error) {
	return pl.reader.Read(p)
}

func (pl *pipeline) Next() *pipeline {
	r, wc := newPipe()
	return &pipeline{
		rw: pl.rw,
		state: &state{
			Named: pl.state.Named,
		},
		// Fresh set of headers to support concurrency
		headers: http.Header{},
		reader:  pl.next,
		writer:  wc,
		next:    r,
	}
}

func newPipe() (io.Reader, io.WriteCloser) {
	pipe := &pipe{
		b: new(bytes.Buffer),
		c: make(chan struct{}),
	}
	return pipe, pipe
}

type pipe struct {
	b *bytes.Buffer
	c chan struct{}
}

func (p *pipe) Write(b []byte) (n int, err error) {
	return p.b.Write(b)
}

func (p *pipe) Close() error {
	close(p.c)
	return nil
}

func (p *pipe) Read(b []byte) (n int, err error) {
	<-p.c
	return p.b.Read(b)
}
