package slot_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/livebud/slot"
	"github.com/matryer/is"
)

func TestChainMainSlot(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
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
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<layout>%s</layout>", slot)
	})
	handler := slot.Batch(view, frame1, frame2, layout)
	// TODO: requests should be able to read the post body. We're being too clever
	// overridding the request body with the pipeline and assuming this will only
	// be used for GET requests.
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"hi"}`))
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
		slots := slot.Open(w, r)
		script := slots.Slot("script")
		err := script.WriteString(`<script src='module.js'></script>`)
		is.NoErr(err)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
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
		slots := slot.Open(w, r)
		w.Header().Set("aa", "aa")
		script := slots.Slot("script")
		err := script.WriteString(`<script src='module.js'></script>`)
		is.NoErr(err)
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", slot)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		w.Header().Set("bb", "bb")
		slot, err := slots.ReadString()
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", slot)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		is.NoErr(err)
		err = slots.Slot("style").WriteString(`<link href='/hi.css'/>`)
		is.NoErr(err)
		w.Header().Set("aa", "cc")
		fmt.Fprintf(w, "<frame2>%s</frame2>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
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

func TestChainOneHandler(t *testing.T) {
	is := is.New(t)
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		is.NoErr(err)
		w.Write(body)
	})
	handler := slot.Chain(h1)
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"hi"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `{"name":"hi"}`)
}

func TestBatchOneHandler(t *testing.T) {
	is := is.New(t)
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		is.NoErr(err)
		w.Write(body)
	})
	handler := slot.Batch(h1)
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"hi"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `{"name":"hi"}`)
}

func ExampleBatch() {
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		w.Header().Set("Content-Type", "text/html")
		script := slots.Slot("script")
		script.WriteString(`<script src='/index.js'></script>`)
		slot, err := slots.ReadString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "<h1>%s</h1>", slot)
	})
	frame := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slots.Slot("style").WriteString(`<link href='/frame.css'/>`)
		fmt.Fprintf(w, "<main>\n\t\t\t%s\n\t\t</main>", slot)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slots := slot.Open(w, r)
		slot, err := slots.ReadString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		script, err := slots.Slot("script").ReadString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		style, err := slots.Slot("style").ReadString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "<html>\n\t<head>\n\t\t%s\n\t\t%s\n\t</head>\n\t<body>\n\t\t%s\n\t</body>\n</html>", script, style, slot)
	})
	handler := slot.Batch(view, frame, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	response, err := httputil.DumpResponse(res, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.ReplaceAll(string(response), "\r\n", "\n"))
	// Output:
	// HTTP/1.1 200 OK
	// Connection: close
	// Content-Type: text/html
	//
	// <html>
	// 	<head>
	// 		<script src='/index.js'></script>
	// 		<link href='/frame.css'/>
	// 	</head>
	// 	<body>
	// 		<main>
	// 			<h1></h1>
	// 		</main>
	// 	</body>
	// </html>
}
