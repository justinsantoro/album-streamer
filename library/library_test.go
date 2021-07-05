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

// root is parsed to build an in memory cache of the library

// the cache will be an in-memory badger db

// key = 0artist/album/song

// 0 prefix is to leave room for using the first byte to denote
// "tables"

//lib.Artists() returns list of artists
//Lib.Albums(Artist) returns an artist's albums
//Lib.Songs(Album) returns a list of songs in an album

const songPath = "Artist1/Album1/song1.mp3"
const artPath = "Artist1/Album1/Album1.webp"
const testDir = "./testdata"

var (
	tArtists = []string{"Artist1", "Artist2"}
	tAlbums  = []string{"Album1", "Album2"}
	tSongs   = []string{"song1.mp3", "song2.mp3", "song3.mp3"}
)

func TestLibrary(t *testing.T) {
	var fsys fs.FS
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
			t.Errorf("bad albums len : %d", len(artists))
			t.FailNow()
		}
		for i, artist := range artists {
			if artist.String() != tArtists[i] {
				t.Errorf("bad artist %s at index %d", artist.String(), i)
				break
			}
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
			t.Errorf("bad albums len: %d", len(albums))
			t.FailNow()
		}
		for i, album := range albums {
			if album.String() != tAlbums[i] {
				t.Errorf("bad album %s at index %d", album.String(), i)
				break
			}
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

		if len(songs) != 2 {
			t.Errorf("bad songs len: %d", len(songs))
		}
		for i, song := range songs {
			if song.String() != tSongs[i] {
				t.Errorf("bad song %s at index %d", song.String(), i)
				break
			}
		}
	})

	t.Run("TestGetSongPath", func(t *testing.T) {
		if path = songs[0].Path(); path != songPath {
			t.Errorf("library gave bad song path: %s - Expected: %s", songs[0].Path(), songPath)
		}
	})
}
