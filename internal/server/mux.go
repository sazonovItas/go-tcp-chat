package tcpws

import (
	"log"

	gotcpws "github.com/sazonovItas/go-tcpws"
)

const (
	ProtoHTTP = "http"
	ProtoWS   = "ws"
)

type (
	// Middleware type for covering handler functions
	Middleware func(next HandlerFunc) HandlerFunc
)

// NewMuxHandler creates new handler for tcpws connection
func NewMuxHandler() *MuxHandler {
	return &MuxHandler{
		middlewares: []Middleware{},
		routerTree: &routingNode{
			children: map[string]*routingNode{},
		},
	}
}

// MuxHandler specifies what handler will handle request
// TODO: replace middlewares with one covered middleware
type MuxHandler struct {
	routerTree *routingNode

	middlewares []Middleware
}

// Add new middleware for request
func (mh *MuxHandler) Use(md Middleware) {
	mh.middlewares = append(mh.middlewares, md)
}

func (mh *MuxHandler) newMiddlewareHandler(h HandlerFunc) HandlerFunc {
	handler := h
	for _, middleware := range mh.middlewares {
		handler = middleware(handler)
	}

	return handler
}

// Set handler function for method and url
// Will panic if can't add handler to routing tree
func (mh *MuxHandler) HandleFunc(method, url string, handler HandlerFunc) {
	p, err := parsePattern(method, url)
	if err != nil {
		panic(err)
	}

	handler = mh.newMiddlewareHandler(handler)

	mh.routerTree.addPattern(p, handler)
}

// Serve connection and call handlers for serving
// TODO: Add logger for serving new connection
func (mh *MuxHandler) Serve(conn *gotcpws.Conn) {
	req, err := conn.ReadFrame()
	if err != nil {
		log.Printf("error to read frame: %s", err.Error())
		return
	}

	request, err := newRequest(req)
	if err != nil {
		log.Printf("error to create new request: %s", err.Error())
		return
	}
	response := newResponse(conn)
	response.Req = request

	n, m := mh.routerTree.match(request.Method, request.Url)
	if n == nil {
		log.Printf("mismatched route")
		return
	}

	// get pattern and matches from url
	request.pattern, request.matches = n.pattern, m

	// TODO: Add default handler for unknown url
	handler := n.handler
	if handler == nil {
		return
	}

	handler.Serve(response, request)

	if request.Proto != ProtoWS {
		_ = response.Write()
	}
}
