package library_test

import (
	"fmt"
	"github.com/justinsantoro/album-streamer/library"
	"io/fs"
	"os"
	"testing"
)

//use a fs.FS
//
//root - /
//     - artist/
//        - album/
//           - song.mp3

//songs will be played in lexical order

// on server start the root is parsed to build an in memory cache of the library

// the cache will be an in-memory badger db

// key = 0artist/album/song
// the album name respectively

// 0 prefix is to leave room for using the first byte to denote
// "tables"

//build the filepath of the song based on the key
//given the root path

//lib.Artists() returns list of artists
//Lib.Albums(Artist) returns an artist's albums
//Lib.Songs(Artist, Album) returns a list of songs in an album

const songPath = "Artist1/Album1/song1.mp3"
const artPath = "Artist1/Album1/Album1.webp"
const testDir = "./testdata"

//const itemlen = 2

func TestLibrary(t *testing.T) {
	var fsys fs.FS
	//TODO: add test FS
	fsys = os.DirFS(testDir)

	var lib *library.Library
	var err error
	lib, err = library.NewLibrary(fsys)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var artists []library.Artist
	t.Run("TestGetArtists", func(t *testing.T) {
		artists, err = lib.Artists()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if len(artists) != 2 {
			t.Errorf("albums len only: %d", len(artists))
			t.FailNow()
		}
	})

	var albums []library.Album
	t.Run("TestGetAlbums", func(t *testing.T) {
		albums, err = lib.Albums(artists[0])
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println(albums)
		if len(albums) != 2 {
			t.Errorf("albums len only: %d", len(albums))
			t.FailNow()
		}
	})

	var path string
	t.Run("TestGetAlbumArtPath", func(t *testing.T) {
		if path = albums[0].Art().Path(); path != artPath {
			t.Errorf("library gave bad artPath: %s - Expected %s", path, artPath)
		}
	})

	var songs []library.Song
	t.Run("TestGetAlbumSongs", func(t *testing.T) {
		songs, err = lib.Songs(albums[0])
		fmt.Println(songs)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("TestGetSongPath", func(t *testing.T) {
		if path = songs[0].Path(); path != songPath {
			t.Errorf("library gave bad songPath: %s - Expected: %s", songs[0].Path(), songPath)
		}
	})
}
