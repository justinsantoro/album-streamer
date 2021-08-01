package h2c

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"net"
	"net/http"
)

type Server struct {
	http2.Server
	BaseConfig *http.Server
	BaseCtx    context.Context
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.BaseConfig.Addr)
	if err != nil {
		return fmt.Errorf("error setting up listener: %v", err)
	}
	if s.BaseCtx == nil {
		s.BaseCtx = context.Background()
	}
	for {
		if err := s.BaseCtx.Err(); err != nil {
			return err
		}
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		}
		s.ServeConn(conn, &http2.ServeConnOpts{
			BaseConfig: s.BaseConfig,
		})
	}
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler http.Handler) error {
	server := &Server{http2.Server{}, &http.Server{Addr: addr, Handler: handler}, nil}
	return server.ListenAndServe()
}

func ListenAndServeWithContext(ctx context.Context, addr string, handler http.Handler) error {
	server := &Server{http2.Server{}, &http.Server{Addr: addr, Handler: handler}, ctx}
	return server.ListenAndServe()
}

var DefaultClient = http.Client{
	Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	},
}
