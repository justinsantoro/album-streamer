package main

import (
	"context"
	"crypto/tls"
	"github.com/hajimehoshi/oto"
	"github.com/justinsantoro/album-streamer/Server/internal"
	"golang.org/x/net/http2"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

var albumToStream = internal.Album{
	Name:      "test",
	Artist:    "test",
	Tracks:    []internal.Track{
		{
			ReaderFunc: func() io.ReadCloser {
				r, err := os.Open("C:\\Users\\jzs\\2_Polygondwanaland.mp3")
				if err != nil {
					panic(err)
				}
				return r
			},
		},
		{
			ReaderFunc: func() io.ReadCloser {
				r, err := os.Open("C:\\Users\\jzs\\1_Crumbling_Castle.mp3")
				if err != nil {
					panic(err)
				}
				return r
			},
		},
	},
}

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
	//log.Println("writing ", len(p), " bytes")
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

	//ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)
	_, err := io.Copy(flushWriter{w: w}, readerCtx{ctx: context.Background(), r: &albumToStream})
	if err != nil {
		log.Println("error streaming album: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		if err != nil {
			log.Println("error accepting connection: ", err)
		}
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

	d, err := mp3.NewDecoder(readerCtx{context.TODO(), resp.Body})
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
