package internal_test

import (
	"bytes"
	"fmt"
	"github.com/justinsantoro/album-streamer/Server/internal"
	"io"
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

var songbytes = append([]byte("americanidiot"), append([]byte("esusofsuburbia"), []byte("olidaysong")...)...)

func TestAlbum_Read(t *testing.T) {
	b := make([]byte, 0)
	for {
		p := make([]byte, 1)
		i, err := a.Read(p)
		fmt.Println(p)
		println(string(p) + " " + fmt.Sprint(i) + " " + fmt.Sprint(len(p)))
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			t.FailNow()
		}
		if i > 0 {
			b = append(b, p[0])
		}
		//println(string(b))
	}

	//b, err := ioutil.ReadAll(&a)
	//if err != nil {
	//	t.Errorf("error reading ablum: %v", err)
	//	t.FailNow()
	//}
	if !bytes.Equal(b, songbytes) {
		t.Errorf("album read returned unexpected bytes: %s - expected %s", string(b), string(songbytes))
		t.FailNow()
	}
}