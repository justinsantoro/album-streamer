module github.com/justinsantoro/album-streamer/streamer

replace github.com/justinsantoro/album-streamer/h2c => ../h2c

replace github.com/justinsantoro/album-streamer/streamer/library => ./library

go 1.16

require (
	github.com/justinsantoro/album-streamer/h2c v0.0.0-00010101000000-000000000000
	github.com/justinsantoro/album-streamer/streamer/library v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
