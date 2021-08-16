module github.com/justinsantoro/album-streamer/player/cmd

replace github.com/justinsantoro/album-streamer/player => ../

replace github.com/justinsantoro/album-streamer/h2c => ../../h2c

go 1.16

require (
	github.com/justinsantoro/album-streamer/h2c v0.0.0-20210816214815-599c40d166a5
	github.com/justinsantoro/album-streamer/player v0.0.0-20210816214815-599c40d166a5
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
)
