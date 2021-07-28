package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/justinsantoro/album-streamer/streamer"
	"log"
	"os"
	"os/signal"
)

func main() {
	artist := flag.String("artist", "", "artist of album to stream")
	album := flag.String("album", "", "name of album to stream")
	lib := flag.String("lib", "", "root dir of library")
	to := flag.String("to", "", "addr of player to stream to")
	flag.Parse()

	if len(*artist) == 0 {
		fmt.Println("forgot --artist")
		os.Exit(1)
	}
	if len(*album) == 0 {
		fmt.Println("forgot --album")
		os.Exit(1)
	}
	if len(*lib) == 0 {
		fmt.Println("forgot --lib")
		os.Exit(1)
	}
	if len(*to) == 0 {
		fmt.Println("forgot --to")
		os.Exit(1)
	}
	strmr, err := streamer.NewStreamer(os.DirFS(*lib))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cfunc := context.WithCancel(context.Background())
	strm, err := strmr.Stream(ctx, *artist, *album, *to)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal, 0)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		cfunc()
	}()
	err = strm.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
