module github.com/justinsantoro/album-streamer/player

replace github.com/justinsantoro/album-streamer/h2c => ../h2c

go 1.16

require (
	github.com/hajimehoshi/go-mp3 v0.3.2
	github.com/hajimehoshi/oto v0.7.2
	github.com/justinsantoro/album-streamer/h2c v0.0.0-00010101000000-000000000000
)
