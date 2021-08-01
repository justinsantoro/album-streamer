package player

import (
	"context"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/justinsantoro/album-streamer/h2c"
	"io"
	"log"
	"net/http"
)

var Handlef http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	log.Println("starting player")
	// First flash response headers
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	//_, err := io.Copy(os.Stdout, r.Body)
	//if err != nil {
	//	log.Println(err)
	//}
	d, err := mp3.NewDecoder(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	d.Length()
	log.Println("starting context")
	//TODO: configure these values via request header
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	log.Println("flashed response")

	if _, err := io.Copy(p, d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("copy end")
}

func ListenAndServe(addr string) error {
	return h2c.ListenAndServe(addr, Handlef)
}

func ListenAndServeWithContext(ctx context.Context, addr string) error {
	return h2c.ListenAndServeWithContext(ctx, addr, Handlef)
}
