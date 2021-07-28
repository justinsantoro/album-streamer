package main

import (
	"flag"
	"github.com/justinsantoro/album-streamer/h2c"
	"github.com/justinsantoro/album-streamer/player"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", "", "listen address")
	flag.Parse()
	serv := h2c.Server{
		Server:     http2.Server{IdleTimeout: time.Second * 3},
		BaseConfig: &http.Server{Addr: *addr, Handler: player.Handlef, WriteTimeout: time.Hour * 1},
		BaseCtx:    nil,
	}
	//log.Fatal(player.ListenAndServe(*addr))
	log.Fatal(serv.ListenAndServe())
}
