package slot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func New(w http.ResponseWriter, r *http.Request) *Slots {
	return &Slots{
		w: w,
		r: r,
	}
}

type once struct {
	once  sync.Once
	state *state
	err   error
}

func (o *once) Do(fn func() (*state, error)) (*state, error) {
	o.once.Do(func() {
		o.state, o.err = fn()
	})
	return o.state, o.err
}

type Slots struct {
	w http.ResponseWriter
	r *http.Request

	once  once
	state *state
}

type state struct {
	Main   string
	Others map[string]string
}

var emptyState = &state{
	Others: map[string]string{},
}

func readState(r io.Reader) (*state, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	} else if len(bytes) == 0 {
		return emptyState, nil
	}
	var state state
	if err := json.Unmarshal(bytes, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *Slots) readState() (*state, error) {
	return s.once.Do(func() (*state, error) {
		return readState(s.r.Body)
	})
}

func (s *Slots) Slot(slot string) *Slot {
	return &Slot{
		w:         s.w,
		name:      slot,
		readState: s.readState,
	}
}

func (s *Slots) ReadString() (string, error) {
	state, err := s.readState()
	if err != nil {
		return "", err
	}
	return state.Main, nil
}

type Slot struct {
	w         http.ResponseWriter
	name      string
	readState func() (*state, error)
}

func (s *Slot) ReadString() (string, error) {
	state, err := s.readState()
	if err != nil {
		return "", err
	}
	return state.Others[s.name], nil
}

type writerTo interface {
	WriteTo(slot string, p []byte) (n int, err error)
}

func (s *Slot) WriteString(str string) error {
	sw, ok := s.w.(writerTo)
	if !ok {
		return fmt.Errorf("slot: responsewriter is not a slot writer")
	}
	_, err := sw.WriteTo(s.name, []byte(str))
	return err
}
