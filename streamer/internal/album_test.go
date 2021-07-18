package internal_test

import (
	"bytes"
	"github.com/justinsantoro/album-streamer/streamer/internal"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

//create album

//get album reader

//streams multiple files via one reader
type stringReadCloser struct {
	io.Reader
}

func (src stringReadCloser) Close() error {
	return nil
}

var a = internal.Album{
	Name: "American Idiot",
	Artist: "Green Day",
	Tracks: []internal.Track{
		internal.Track{
			Name: "American Idiot",
			ReaderFunc: func () io.ReadCloser {
				return stringReadCloser{strings.NewReader("americanidiot")}
			},
		},
		internal.Track{
			Name: "Jesus Of Suburbia",
			ReaderFunc: func () io.ReadCloser {
				return stringReadCloser{strings.NewReader("jesusofsuburbia")}
			},
		},
		internal.Track{
			Name: "Holiday",
			ReaderFunc: func() io.ReadCloser {
				return stringReadCloser{strings.NewReader("holidaysong")}
			},
		},
	},
}

var songbytes = append([]byte("americanidiot"), append([]byte("jesusofsuburbia"), []byte("holidaysong")...)...)

func TestAlbum_Read(t *testing.T) {
	b, err := ioutil.ReadAll(&a)
	if err != nil {
		t.Errorf("error reading ablum: %v", err)
		t.FailNow()
	}
	if !bytes.Equal(b, songbytes) {
		t.Errorf("album read returned unexpected bytes: %s - expected %s", string(b), string(songbytes))
		t.FailNow()
	}
}