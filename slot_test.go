package slot_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/slot"
	"github.com/matryer/is"
)

func TestChainMainSlot(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<layout>%s</layout>", slot)
	})
	handler := slot.Chain(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestBatchMainSlot(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<layout>%s</layout>", slot)
	})
	handler := slot.Batch(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestChainOtherSlots(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		script := slots.Slot("script")
		err := script.WriteString(`<script src='module.js'></script>`)
		is.NoErr(err)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		script, err := slots.Slot("script").ReadString()
		is.NoErr(err)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<layout>%s%s</layout>", script, slot)
	})
	handler := slot.Chain(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><script src='module.js'></script><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestBatchOtherSlots(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		w.Header().Set("aa", "aa")
		script := slots.Slot("script")
		err := script.WriteString(`<script src='module.js'></script>`)
		is.NoErr(err)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		w.Header().Set("bb", "bb")
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		err = slots.Slot("style").WriteString(`<link href='/hi.css'/>`)
		is.NoErr(err)
		w.Header().Set("aa", "cc")
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.New(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		script, err := slots.Slot("script").ReadString()
		is.NoErr(err)
		style, err := slots.Slot("style").ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<layout>%s%s%s</layout>", script, style, slot)
	})
	handler := slot.Batch(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `<layout><script src='module.js'></script><link href='/hi.css'/><frame2><frame1><view></view></frame1></frame2></layout>`)
	headers := res.Header
	is.Equal(headers.Get("aa"), "aa")
	is.Equal(headers.Get("bb"), "bb")
}
