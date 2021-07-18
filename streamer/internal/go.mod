module github.com/justinsantoro/album-streamer/Server/internal

replace github.com/justinsantoro/album-streamer/h2c => ../../h2c

go 1.16

require (
	github.com/justinsantoro/album-streamer/h2c v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
