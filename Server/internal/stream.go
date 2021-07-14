package internal

import (
	"context"
	"fmt"
	"github.com/justinsantoro/album-streamer/h2c"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"net/http"
)


type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

func (r readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, io.EOF
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

type Stream struct {
	c http.Client
	ctx context.Context
	cfunc context.CancelFunc
	a *Album
	to string
	w func() error
}

//NewStream returns starns and returns an new stream to the given url
func NewStream(ctx context.Context, album *Album, to string) (*Stream, error) {
	var strm Stream
	ctx, cfunc := context.WithCancel(ctx)
	strm = Stream{
		c:   h2c.DefaultClient,
		ctx: ctx,
		cfunc: cfunc,
		a:   album,
		to: to,
	}
	if err := strm.start() ;err != nil {
		return nil, err
	}
	return &strm, nil
}

func (s *Stream) start() error {
	pr, pw := io.Pipe()

	req, err := http.NewRequestWithContext(s.ctx, http.MethodPut, s.to, io.NopCloser(pr))
	if err != nil {
		return err
	}

	// Send the request
	resp, err := s.c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s reponded with unssuccessful status code: %d", s.to, resp.StatusCode)
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		_, err := io.Copy(pw, readerCtx{ctx: s.ctx, r: s.a})
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		//TODO: do something with the incoming data stream?
		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil && err != context.Canceled {
			return err
		}
		return nil
	})
	s.w = eg.Wait
	return nil
}

func (s *Stream) Wait() error {
	return s.w()
}

