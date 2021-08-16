module github.com/justinsantoro/streamer/cmd

replace github.com/justinsantoro/album-streamer/streamer => ../

replace github.com/justinsantoro/album-streamer/streamer/library => ../library

replace github.com/justinsantoro/album-streamer/streamer/internal => ../internal

replace github.com/justinsantoro/album-streamer/h2c => ../../h2c

go 1.16

require github.com/justinsantoro/album-streamer/streamer v0.0.0-20210816214815-599c40d166a5
