package main

import (
	"flag"
	"github.com/justinsantoro/album-streamer/player"
	"log"
)

func main() {
	addr := flag.String("addr", "", "listen address")
	flag.Parse()
	log.Fatal(player.ListenAndServe(*addr))
}
