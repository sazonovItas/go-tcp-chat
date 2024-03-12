package tcpws

import gotcpws "github.com/sazonovItas/go-tcpws"

const (
	ProtoHTTP = "HTTP"
	ProtoWS   = "WS"
)

// NewMuxHandler creates new handler for tcpws connection
func NewMuxHandler() *MuxHandler {
	return &MuxHandler{
		routerTree: &routingNode{
			children: map[string]*routingNode{},
		},
	}
}

// MuxHandler specifies what handler will handle request
type MuxHandler struct {
	routerTree *routingNode
}

// Set handler function for method and url
// Will panic if something is wrong
func (mh *MuxHandler) HandleFunc(method, url string, handler HandlerFunc) {
	p, err := parsePattern(method, url)
	if err != nil {
		panic(err)
	}

	mh.routerTree.addPattern(p, handler)
}

// Serve connection and call handlers for serving
func (mh *MuxHandler) Serve(conn *gotcpws.Conn) {
	req, err := conn.ReadFrame()
	if err != nil {
		return
	}

	request, err := newRequest(req)
	if err != nil {
		return
	}
	response := newResponse(conn)
	response.Req = request

	n, m := mh.routerTree.match(request.Method, request.Url)
	if n == nil {
		return
	}

	// get pattern and matches from url
	request.pattern, request.matches = n.pattern, m

	handler := n.handler
	if handler == nil {
		return
	}

	handler.Serve(response, request)
}
