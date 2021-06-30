package library

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
)

const (
	webpExt = ".webp"
	pathSep = "/"
	)
const (
	iartist = iota
	ialbum
	isong
)
type libkey []byte

func (lk libkey) String() string {
	return string(lk[1:])
}

func parseKey(k libkey, i int) string {
	return string(bytes.Split(k[1:], []byte(pathSep))[i])
}

// Artist is the name of an artist by which
// Albums are contained in the library
type Artist struct {
	libkey
}

// String returns the Artist name as a string
func (a Artist) String() string {
	return parseKey(a.libkey, iartist)
}

//Path returns the path to the artist
func (a Artist) Path() string {
	return fmt.Sprintf("/%s", a)
}

// Album is an album of music contained in the
// library
type Album struct {
	libkey
}

// String retures the string of the Album name
func (a Album) String() string {
	return parseKey(a.libkey, ialbum)
}

//Path returns the path to the Album
func (a Album) Path() string {
	return fmt.Sprintf("/%s/%s/", a.Artist(), a)
}

// Artist returns the Artist which the album is
// by
func (a Album) Artist() Artist {
	return Artist{a.libkey}
}

// AlbumArt is the artwork associated with an
// album. The artwork must be a webp image which
// bears the same name as the album + webpExt
type AlbumArt struct {
	album Album
}

//Path returns the path to the AlbumArt
func (a AlbumArt) Path() string {
	return a.album.Path() + a.album.String() + webpExt
}

// Reader returns a ReadCloser for the AlbumArt webp file
func (a AlbumArt) Reader(fsys fs.FS) (io.ReadCloser, error) {
	return fsys.Open(a.Path())
}

// Art returns the AlbumArt associated with the Album
func (a Album) Art() AlbumArt {
	return AlbumArt{a}
}

// Song is a song contained in the library
type Song struct {
	libkey
}

// String returns the song name as a string
func (s Song) String() string {
	return parseKey(s.libkey, isong)
}

// Artist returns the Artist which the song is by
func (s Song) Artist() Artist {
	return Artist{s.libkey}
}

// Album returns the Album the song is on
func (s Song) Album() Album {
	return Album{s.libkey}
}

// Path returns the path to the Album
func (s Song) Path() string {
	return fmt.Sprintf("/%s/%s/%s", s.Artist(), s.Album(), s)
}

// Reader returns a ReadCloser for the song file
func (s Song) Reader(fsys fs.FS) (io.ReadCloser, error) {
	return fsys.Open(s.Path())
}