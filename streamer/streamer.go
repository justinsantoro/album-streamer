package streamer

import (
	"context"
	"fmt"
	"github.com/justinsantoro/album-streamer/streamer/internal"
	"github.com/justinsantoro/album-streamer/streamer/library"
	"io"
	"io/fs"
)

type Streamer struct {
	fsys fs.FS
	*library.Library
}

func NewStreamer(fsys fs.FS) (*Streamer, error) {
	l, err := library.NewLibrary(fsys)
	if err != nil {
		return nil, err
	}
	return &Streamer{fsys, l}, nil
}

func (s *Streamer) Stream(ctx context.Context, artist string, album string, to string) (*internal.Stream, error) {
	a := library.Album(fmt.Sprint(artist, "/", album))
	sngs, err := s.Library.Songs(a)
	if err != nil {
		return nil, fmt.Errorf("error getting album %s's songs: %v", a, err)
	}
	trks := make([]internal.Track, 0)
	for _, trk := range sngs {
		t := trk
		trks = append(trks, internal.Track{
			Name: t.String(),
			ReaderFunc: func() (io.ReadCloser, error) {
				return t.Reader(s.fsys)
			},
		})
	}
	alb := &internal.Album{
		Name:   a.String(),
		Artist: a.Artist().String(),
		Tracks: trks,
		ArtReader: func() (io.ReadCloser, error) {
			return a.Art().Reader(s.fsys)
		},
	}
	return internal.NewStream(ctx, alb, to)
}
