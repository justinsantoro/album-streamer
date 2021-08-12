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
	Name:   "American Idiot",
	Artist: "Green Day",
	Tracks: []internal.Track{
		internal.Track{
			Name: "American Idiot",
			ReaderFunc: func() (io.ReadCloser, error) {
				return stringReadCloser{strings.NewReader("americanidiot")}, nil
			},
		},
		internal.Track{
			Name: "Jesus Of Suburbia",
			ReaderFunc: func() (io.ReadCloser, error) {
				return stringReadCloser{strings.NewReader("jesusofsuburbia")}, nil
			},
		},
		internal.Track{
			Name: "Holiday",
			ReaderFunc: func() (io.ReadCloser, error) {
				return stringReadCloser{strings.NewReader("holidaysong")}, nil
			},
		},
	},
}

//for now, album reader always skips at least the first 3 bytes of each songs mp3 file
var songbytes = append([]byte("ricanidiot"), append([]byte("usofsuburbia"), []byte("idaysong")...)...)
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
