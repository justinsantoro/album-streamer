package internal_test

import (
	"context"
	"github.com/justinsantoro/album-streamer/h2c"
	"github.com/justinsantoro/album-streamer/streamer/internal"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type ForeverReader byte

func (fr ForeverReader) Read(p []byte) (int, error) {
	return 1, nil
}

var album = internal.Album{
	Name:   "test",
	Artist: "test artist",
	Tracks: []internal.Track{
		internal.Track{
			Name: "test track",
			ReaderFunc: func() (io.ReadCloser, error) {
				return io.NopCloser(new(ForeverReader)), nil
			},
		},
	},
	ArtReader: nil,
}

func TestServer_Stream(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// First flash response headers
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		_, _ = io.Copy(os.Stdout, r.Body)
	})

	go h2c.ListenAndServe("0.0.0.0:8080", nil)

	var err error
	ctx, cfunc := context.WithCancel(context.Background())
	wait, err := internal.NewStream(ctx, &album, "http://0.0.0.0:8080")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ch := make(chan error, 0)
	go func() {
		ch <- wait()
	}()

	cfunc()

	select {
	case <-time.After(time.Millisecond * 500):
		t.Errorf("strm.Wait timeout")
		t.FailNow()
	case err = <-ch:
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
