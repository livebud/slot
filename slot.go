package slot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Open slots for reading and writing
func Open(w http.ResponseWriter, r *http.Request) *slots {
	return &slots{
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

// NamedSlot is a single read-writable slot
type NamedSlot interface {
	ReadString() (string, error)
	WriteString(string) error
}

// Slots contains a single readable default slot and a map of named slots
type Slots interface {
	ReadString() (string, error)
	Slot(slot string) NamedSlot
}

type slots struct {
	w    http.ResponseWriter
	r    *http.Request
	once once
}

var _ Slots = (*slots)(nil)

// State of the slots. This data structure may change, so don't rely on this
// structure in your app.
type state struct {
	Data  string            `json:"data"`
	Named map[string]string `json:"named"`
}

var emptyState = &state{
	Named: map[string]string{},
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

func (s *slots) readState() (*state, error) {
	return s.once.Do(func() (*state, error) {
		return readState(s.r.Body)
	})
}

func (s *slots) Slot(name string) NamedSlot {
	return &namedSlot{
		w:         s.w,
		name:      name,
		readState: s.readState,
	}
}

func (s *slots) ReadString() (string, error) {
	state, err := s.readState()
	if err != nil {
		return "", err
	}
	return state.Data, nil
}

type namedSlot struct {
	w         http.ResponseWriter
	name      string
	readState func() (*state, error)
}

var _ NamedSlot = (*namedSlot)(nil)

func (s *namedSlot) ReadString() (string, error) {
	state, err := s.readState()
	if err != nil {
		return "", err
	}
	return state.Named[s.name], nil
}

type writerTo interface {
	WriteTo(slot string, p []byte) (n int, err error)
}

func (s *namedSlot) WriteString(str string) error {
	sw, ok := s.w.(writerTo)
	if !ok {
		return fmt.Errorf("slot: responsewriter is not a slot writer")
	}
	_, err := sw.WriteTo(s.name, []byte(str))
	return err
}
