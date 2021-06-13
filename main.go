package main

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/http2"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"github.com/hajimehoshi/oto"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

const mp3fpath = "C:\\Users\\jzs\\classic.mp3"

type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

func (r readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	log.Println("writing ", len(p), " bytes")
	// Flush - send the buffered written data to the client
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

func streamAudio(w http.ResponseWriter, r *http.Request) {
	// First flash response headers
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	f, err := os.Open(mp3fpath)
	if err != nil {
		log.Println("error opening file to stream: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//need reader in between here that only sends eof at end of album not at end of files
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)
	_, _ = io.Copy(flushWriter{w: w}, readerCtx{ctx: ctx, r: f})

}

func H2CServerPrior() {
	server := http2.Server{}

	l, err := net.Listen("tcp", "0.0.0.0:1010")
	if err != nil {
		log.Fatal("error setting up listener: ", err)
	}

	log.Println("Listening [0.0.0.0:1010]...")
	for {
		conn, err := l.Accept()
		log.Println("error accepting connection: ", err)

		server.ServeConn(conn, &http2.ServeConnOpts{
			Handler: http.HandlerFunc(streamAudio),
		})
	}
}

func GetAudioStream(ctx context.Context) error {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	resp, err := client.Get("http://0.0.0.0:1010")
	if err != nil {
		return err
	}

	d, err := mp3.NewDecoder(readerCtx{ctx, resp.Body})
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()


	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil

}

func main() {
	log.SetOutput(os.Stdout)
	//start server
	go H2CServerPrior()
	//stream audio
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 8)
	log.Fatal(GetAudioStream(ctx))
}
