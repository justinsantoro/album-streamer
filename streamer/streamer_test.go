package streamer_test

import (
	"context"
	"github.com/justinsantoro/album-streamer/h2c"
	"github.com/justinsantoro/album-streamer/streamer"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

const expectedStream = "sissong1sissong2"
const testDir = "./library/testdata"

func TestStream(t *testing.T) {
	var fsys fs.FS
	fsys = os.DirFS(testDir)

	ch := make(chan interface{}, 0)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// First flash response headers
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ch <- err
			return
		}
		ch <- string(b)
	})

	go h2c.ListenAndServe("0.0.0.0:8080", nil)

	strmr, err := streamer.NewStreamer(fsys)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	wait, err := strmr.Stream(context.Background(), "Artist1", "Album1", "http://0.0.0.0:8080")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	go func() {
		ch <- wait()
	}()

	select {
	case <-time.After(time.Millisecond * 500):
		t.Errorf("strm.Wait timeout")
		t.FailNow()
	case k := <-ch:
		switch v := k.(type) {
		case error:
			t.Error(v)
			t.FailNow()
		case string:
			if v != expectedStream {
				t.Errorf("unexpected stream output: %s - expected %s", v, expectedStream)
				t.FailNow()
			}
		}
	}
}
