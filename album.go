package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type track struct {
	fpath string
	name string
}

type album struct {
	name string
	artist string
	tracks []track
	art string
	cTrack io.ReadCloser
	cTrackNum int
}

func (a *album) currentTrack() track {
	return a.tracks[a.cTrackNum]
}

func (a *album) nextTrack() error {
	if a.cTrack != nil {
		a.cTrack.Close()
		a.cTrackNum++
	}
	if a.cTrackNum + 1 > len(a.tracks) {
		//end of album
		return io.EOF
	}
	var err error
	a.cTrack, err = os.Open(a.currentTrack().fpath)
	if err != nil {
		return fmt.Errorf("error opening track %s: %v", a.currentTrack().name, err)
	}
	log.Printf("playing track %d", a.cTrackNum)
	return nil
}

func (a *album) Read(p []byte) (int, error) {
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
		//fill p with bytes from next track
		b := make([]byte, len(p))
		j, err := a.cTrack.Read(b)
		if err != nil {
			return j, err
		}
		p = b
	}
	return i, nil
}