package library

import (
	"fmt"
	"io"
	"io/fs"
)

const webpExt = ".webp"

type libkey []byte

func (lk libkey) artistLen() int8 {
	return int8(lk[len(lk)-3])
}

func (lk libkey) albumLen() int8 {
	return int8(lk[len(lk)-2])
}

// Artist is the name of an artist by which
// Albums are contained in the library
type Artist struct {
	libkey
}

// Bytes returns the bytes of the Artist name
func (a Artist) Bytes() []byte {
	return a.libkey[1:a.libkey.artistLen()-1]
}

// String returns the Artist name as a string
func (a Artist) String() string {
	return string(a.Bytes())
}

// Album is an album of music contained in the
// library
type Album struct {
	libkey
}

// Bytes returns the bytes of the Album name
func (a Album) Bytes() []byte {
	return a.libkey[1:a.libkey.albumLen()-1]
}

// String retures the string of the Album name
func (a Album) String() string {
	return string(a.Bytes())
}

func (a Album) Path() string {
	return fmt.Sprintf("/%s/%s/")
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

// Bytes returns the song name as bytes
func (s Song) Bytes() []byte {
	return s.libkey[1:len(s.libkey) -4]
}

// String returns the song name as a string
func (s Song) String() string {
	return string(s.Bytes())
}

// Artist returns the Artist which the song is by
func (s Song) Artist() Artist {
	return Artist{s.libkey}
}

// Album returns the Album the song is on
func (s Song) Album() Album {
	return Album{s.libkey}
}

func (s Song) Path() string {
	return fmt.Sprintf("/%s/%s/%s", s.Artist(), s.Album(), s)
}

// Reader returns a ReadCloser for the song file
func (s Song) Reader(fsys fs.FS) (io.ReadCloser, error) {
	return fsys.Open(s.Path())
}