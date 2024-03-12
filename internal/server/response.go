package tcpws

import (
	"encoding/json"

	gotcpws "github.com/sazonovItas/go-tcpws"
)

// Response for the
type Response struct {
	Header map[string]interface{}
	Body   string

	Req *Request

	// Connection using for ws like communication
	Conn *gotcpws.Conn
}

func newResponse(conn *gotcpws.Conn) *Response {
	return &Response{
		Conn:   conn,
		Header: map[string]interface{}{},
	}
}

func (resp *Response) Write() error {
	type response struct {
		Header map[string]interface{} `json:"header"`
		Body   string                 `json:"body"`
	}

	wrtResp := response{
		Header: resp.Header,
		Body:   resp.Body,
	}

	err := json.NewEncoder(resp.Conn).Encode(wrtResp)
	return err
}
