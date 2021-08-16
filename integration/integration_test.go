package main

import (
	"context"
	"github.com/justinsantoro/album-streamer/player"
	"github.com/justinsantoro/album-streamer/streamer"
	"os"
	"testing"
)

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
		wait, err := s.Stream(context.Background(), "Artist1", "Album1", "http://"+paddr)
		if err != nil {
			t.Error(err)
		}
		if err := wait(); err != nil {
			t.Error(err)
		}
	})
	//time.Sleep(time.Second * 5)
}
