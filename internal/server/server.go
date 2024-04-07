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
	}
	return server.ListenAndServe()
}

// NewServer creates new server with handler
func NewServer(addr string, handler HandleFunc) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
		connwg:  &sync.WaitGroup{},
	}
}

// Server is struct for accepting and serving connections
type Server struct {
	Addr    string
	Handler HandleFunc

	// ln is listener for addr
	ln net.Listener

	// connwg wait group for waiting until all serving connections are done
	connwg *sync.WaitGroup
}

// ListenAndServe create listener on server addr, accepting
// connections and serve connection in goroutine
// TODO: remove log from the accpeting connection
func (srv *Server) ListenAndServe() error {
	// create new listener on addr
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer func() {
		srv.connwg.Wait()

		listener.Close()
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
				srv.connwg.Done()
				conn.Close()

				log.Println("closed connection:", conn.RemoteAddr())
			}()

			srv.connwg.Add(1)
			srv.Serve(conn)
		}()
	}
}

// Serve serves connection with server handler
// if server handler is not setuped will panic
func (srv *Server) Serve(conn *gotcpws.Conn) {
	if srv.Handler == nil {
		panic("server handler is not set")
	}

	srv.Handler.Serve(conn)
}
