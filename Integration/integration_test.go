package main

import (
	"context"
	"github.com/justinsantoro/album-streamer/player"
	"github.com/justinsantoro/album-streamer/streamer"
	"io"
	"io/fs"
	"os"
	"testing"
)

type testreader struct {
	io.ReadCloser
	fs *testfs
}

type testfs struct {
	fs.FS
	t testing.T
	ri int
}

func (tr testreader) Read(p []byte) (int, error) {
	i, err := tr.Read(p)
	tr.fs.ri += i
	return i, err
}

func (tfs *testfs) Open(name string) (io.ReadCloser, error) {
	r, err := tfs.FS.Open(name)
	if err != nil {
		tfs.t.Errorf("testfs.Open err: %v", err)
		tfs.t.FailNow()
	}
	return testreader{r, tfs}, err
}

func TestIntegration(t *testing.T) {
	paddr := "127.0.0.1:8325"
	t.Run("NewPlayer", func(t *testing.T) {
		t.Log("starting player at: ", paddr)
		go func() {
			err := player.ListenAndServe(paddr)
			if err != nil {
				t.Log(err)
			}
		}()
	})

	var s *streamer.Streamer
	t.Run("NewStreamer", func(t *testing.T) {
		var err error
		s, err = streamer.NewStreamer(os.DirFS("./testdata"))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("StreamAlbum", func(t *testing.T) {
		strm, err := s.Stream(context.Background(), "Artist1", "Album1", "http://" + paddr)
		if err != nil {
			t.Error(err)
		}
		if err := strm.Wait(); err != nil {
			t.Error(err)
		}
	})
	//time.Sleep(time.Second * 5)
}
