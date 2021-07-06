package internal

import (
	"io"
	"log"
)

type Track struct {
	Name string
	r io.ReadCloser
	ReaderFunc func() io.ReadCloser
}

func (t *Track) Read(p []byte) (int, error) {
	if t.r == nil {
		t.r = t.ReaderFunc()
	}
	return t.r.Read(p)
}

func (t *Track) Close() error {
	if t.r == nil {
		return nil
	}
	return t.r.Close()
}

type Album struct {
	Name string
	Artist string
	Tracks []Track
	Art io.ReadCloser
	cTrack *Track
	cTrackNum int

}

func (a *Album) currentTrack() *Track {
	return &a.Tracks[a.cTrackNum]
}

func (a *Album) nextTrack() error {
	if a.cTrack != nil {
		a.cTrack.Close()
		a.cTrackNum++
	}
	if a.cTrackNum + 1 > len(a.Tracks) {
		//end of album
		return io.EOF
	}
	a.cTrack = a.currentTrack()
	log.Printf("playing track %d", a.cTrackNum)
	return nil
}

func (a *Album) Read(p []byte) (int, error) {
	if a.cTrack == nil {
		if err := a.nextTrack(); err != nil {
			return 0, err
		}
	}
	i, err := a.cTrack.Read(p)
	if err == io.EOF {
		err = a.nextTrack()
		if err != nil {
			return i, err
		}
		return 0, nil
	}
	//println(string(p) + " " + fmt.Sprint(i))
	return i, nil
}