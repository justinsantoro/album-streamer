module github.com/justinsantoro/album-streamer/integration

replace github.com/justinsantoro/album-streamer/h2c => ../h2c

replace github.com/justinsantoro/album-streamer/streamer => ../streamer

replace github.com/justinsantoro/album-streamer/player => ../player

replace github.com/justinsantoro/album-streamer/streamer/library => ../streamer/library

go 1.16

require (
	github.com/justinsantoro/album-streamer/player v0.0.0-00010101000000-000000000000
	github.com/justinsantoro/album-streamer/streamer v0.0.0-00010101000000-000000000000
)
