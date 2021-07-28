package library

import (
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/justinsantoro/wrappedbadger"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
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

type libkey string

func (k libkey) String() string {
	return string(k)
}

func parseKey(k libkey, i int) string {
	return string(strings.Split(string(k), pathSep)[i])
}

func (k libkey) Bytes() []byte {
	return append([]byte{0}, []byte(k)...)
}

func libkeyFromBytes(k []byte) libkey {
	return libkey(string(k[1:]))
}

// Artist is the name of an artist by which
// Albums are contained in the library
type Artist libkey

// String returns the Artist name as a string
func (a Artist) String() string {
	return parseKey(libkey(a), iartist)
}

//Path returns the path to the artist
func (a Artist) Path() string {
	return fmt.Sprintf("%s/", a)
}

// Album is an album of music contained in the
// library
type Album libkey

// String retures the string of the Album name
func (a Album) String() string {
	return parseKey(libkey(a), ialbum)
}

//Path returns the path to the Album
func (a Album) Path() string {
	return fmt.Sprintf("%s/%s/", a.Artist(), a)
}

// Artist returns the Artist which the album is
// by
func (a Album) Artist() Artist {
	return Artist(a)
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
type Song libkey

// String returns the song name as a string
func (s Song) String() string {
	return parseKey(libkey(s), isong)
}

// Artist returns the Artist which the song is by
func (s Song) Artist() Artist {
	return Artist(s)
}

// Album returns the Album the song is on
func (s Song) Album() Album {
	return Album(s)
}

// Path returns the path to the Album
func (s Song) Path() string {
	return fmt.Sprintf("%s/%s/%s", s.Artist(), s.Album(), s)
}

// Reader returns a ReadCloser for the song file
func (s Song) Reader(fsys fs.FS) (io.ReadCloser, error) {
	return fsys.Open(s.Path())
}

type Library struct {
	fsys  fs.FS
	store *wrappedbadger.Store
}

func (l *Library) parse() error {
	return fs.WalkDir(l.fsys, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if filepath.Ext(path) != webpExt {
				fmt.Println(path)
				return l.store.Set(libkey(path).Bytes(), nil)
			}
		}
		return nil
	})
}

func NewLibrary(fsys fs.FS) (*Library, error) {
	var lib Library
	lib.fsys = fsys

	opts := badger.DefaultOptions("")
	opts.Logger = nil
	opts.InMemory = true

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	lib.store = &wrappedbadger.Store{db}
	if err = lib.parse(); err != nil {
		return nil, err
	}
	return &lib, nil
}

func (l Library) Artists() ([]Artist, error) {
	var prevArtist Artist
	artists := make([]Artist, 0)

	err := l.store.IterateKeys(libkey("").Bytes(), func(k []byte) error {
		artist := Artist(libkeyFromBytes(k))
		if artist.String() != prevArtist.String() {
			prevArtist = artist
			artists = append(artists, artist)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func (l Library) Albums(artist Artist) ([]Album, error) {
	prevAlbum := ""
	albums := make([]Album, 0)
	fmt.Println(artist.Path())
	err := l.store.IterateKeys(libkey(artist.Path()).Bytes(), func(k []byte) error {
		album := Album(libkeyFromBytes(k))
		fmt.Println("album key: " + libkey(album).String())
		fmt.Println("album: " + album.String())
		fmt.Println("prev album: " + prevAlbum)
		if album.String() != prevAlbum {
			fmt.Println("set prev album: from " + prevAlbum)
			prevAlbum = album.String()
			fmt.Println("set prev album: to " + prevAlbum)
			albums = append(albums, album)
		}
		fmt.Println("end check album")
		fmt.Println()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (l *Library) Songs(album Album) ([]Song, error) {
	songs := make([]Song, 0)
	fmt.Println(album.Path())
	err := l.store.IterateKeys(libkey(album.Path()).Bytes(), func(k []byte) error {
		songs = append(songs, Song(libkeyFromBytes(k)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return songs, nil
}
 func (l Library) Album(path string) (Album, error) {
 	b, err := l.store.Get([]byte(path))
 	if err != nil {
 		return "", err
	}
	return Album(libkey(b)), nil
 }