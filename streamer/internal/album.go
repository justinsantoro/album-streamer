package internal

import (
	"io"
	"log"
)

type Track struct {
	Name       string
	r          io.ReadCloser
	ReaderFunc func() (io.ReadCloser, error)
}

func (t *Track) Read(p []byte) (int, error) {
	var err error
	if t.r == nil {
		t.r, err = t.ReaderFunc()
		if err != nil {
			return 0, err
		}
		s := source{reader: t.r}
		err = s.skipTags()
		return 0, err
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
	Name      string
	Artist    string
	Tracks    []Track
	cTrack    *Track
	cTrackNum int
	ArtReader func() (io.ReadCloser, error)
}

func (a *Album) currentTrack() *Track {
	return &a.Tracks[a.cTrackNum]
}

func (a *Album) nextTrack() error {
	if a.cTrack != nil {
		a.cTrack.Close()
		a.cTrackNum++
	}
	if a.cTrackNum+1 > len(a.Tracks) {
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
