package main

import (
	"fmt"
	"io"
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
	a.cTrack.Close()
	a.cTrackNum++
	if a.cTrackNum + 1 > len(a.tracks) {
		//end of album
		return io.EOF
	}
	var err error
	a.cTrack, err = os.Open(a.currentTrack().fpath)
	if err != nil {
		return fmt.Errorf("error opening track %s: %v", a.currentTrack().name, err)
	}
	return nil
}

func (a *album) Read(p []byte) (int, error) {
	i, err := a.cTrack.Read(p)
	if err == io.EOF {
		return i, a.nextTrack()
	}
	return i, nil
}