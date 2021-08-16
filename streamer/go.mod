module github.com/justinsantoro/album-streamer/streamer

replace github.com/justinsantoro/album-streamer/h2c => ../h2c

replace github.com/justinsantoro/album-streamer/streamer/library => ./library

go 1.16

require (
	github.com/justinsantoro/album-streamer/h2c v0.0.0-20210816214815-599c40d166a5
	github.com/justinsantoro/album-streamer/streamer/library v0.0.0-20210816214815-599c40d166a5
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
