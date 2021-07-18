package player

import (
	"context"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/justinsantoro/album-streamer/h2c"
	"io"
	"net/http"
)

var handlef http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {

	d, err := mp3.NewDecoder(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: configure these values via request header
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer c.Close()


	p := c.NewPlayer()
	defer p.Close()

	// First flash response headers
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	if _, err := io.Copy(p, d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ListenAndServe(addr string) error {
	return h2c.ListenAndServe(addr, handlef)
}

func ListenAndServeWithContext(ctx context.Context, addr string) error {
	return h2c.ListenAndServeWithContext(ctx, addr, handlef)
}
