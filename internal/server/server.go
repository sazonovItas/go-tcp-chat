package tcpws

import (
	"errors"
	"log"
	"net"
	"sync"

	gotcpws "github.com/sazonovItas/go-tcpws"
)

type HandlerFunc func(resp *Response, req *Request)

func (hf HandlerFunc) Serve(resp *Response, req *Request) {
	hf(resp, req)
}

// HandleFunc is interface for handle gotcpws connection
type HandleFunc interface {
	Serve(conn *gotcpws.Conn)
}

// ListenAndServe creates new server
func ListenAndServe(addr string, handler HandleFunc) error {
	server := &Server{
		Addr:    addr,
		Handler: handler,
		connwg:  &sync.WaitGroup{},
		conns:   map[*gotcpws.Conn]struct{}{},
	}
	return server.ListenAndServe()
}

func NewServer(addr string, handler HandleFunc) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
		connwg:  &sync.WaitGroup{},
		conns:   map[*gotcpws.Conn]struct{}{},
	}
}

// Server is struct for accepting and serving connections
type Server struct {
	Addr    string
	Handler HandleFunc

	// ln is listener for addr
	ln net.Listener

	// conns is storage for active connections
	conns map[*gotcpws.Conn]struct{}

	// connwg wait group for waiting until all serving connections are done
	connwg *sync.WaitGroup
}

// ListenAndServe create listener on server addr, accepting
// connections and serve connection in goroutine
func (srv *Server) ListenAndServe() error {
	// create new listener on addr
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer func() {
		listener.Close()
		srv.connwg.Wait()
	}()

	srv.ln = listener
	for {
		// Accept new connection if accuse error then check error
		// if it isn't ErrClosed continue accepting other connections
		c, err := srv.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			continue
		}

		// create new ws like connection for handling
		conn := gotcpws.NewFrameConnection(c, nil, nil, 0, true)

		log.Println("accepted new connection:", conn.RemoteAddr())

		// Serve connection
		go func() {
			defer func() {
				delete(srv.conns, conn)
				srv.connwg.Done()
				log.Println("closed connection:", conn.RemoteAddr())
				conn.Close()
			}()

			srv.conns[conn] = struct{}{}
			srv.connwg.Add(1)
			srv.Serve(conn)
		}()
	}
}

func (srv *Server) Serve(conn *gotcpws.Conn) {
	if srv.Handler == nil {
		panic("server handler is not set")
	}

	srv.Handler.Serve(conn)
}
