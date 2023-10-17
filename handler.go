package slot

import (
	"net/http"

	"golang.org/x/sync/errgroup"
)

func Chain(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pipeline := newPipeline(w)
		for i := 0; i < len(handlers); i++ {
			handler := handlers[i]
			r := r.Clone(r.Context())
			r.Body = pipeline
			handler.ServeHTTP(pipeline, r)
			pipeline.Close()
			pipeline = pipeline.Next()
		}
		state, err := readState(pipeline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(state.Main))
	})
}

func Batch(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pipeline := newPipeline(w)
		eg, ctx := errgroup.WithContext(r.Context())
		for i := 0; i < len(handlers); i++ {
			handler := handlers[i]
			innerPipeline := pipeline
			r := r.Clone(ctx)
			r.Body = innerPipeline
			eg.Go(func() (err error) {
				defer func() { err = innerPipeline.Close() }()
				handler.ServeHTTP(innerPipeline, r)
				return err
			})
			pipeline = pipeline.Next()
		}
		if err := eg.Wait(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		state, err := readState(pipeline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(state.Main))
	})
}
