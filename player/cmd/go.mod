module github.com/justinsantoro/album-streamer/player/cmd

replace github.com/justinsantoro/album-streamer/player => ../

replace github.com/justinsantoro/album-streamer/h2c => ../../h2c

go 1.16

require (
	github.com/justinsantoro/album-streamer/h2c v0.0.0-00010101000000-000000000000
	github.com/justinsantoro/album-streamer/player v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
)
