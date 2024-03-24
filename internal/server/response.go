package tcpws

import (
	"encoding/json"

	gotcpws "github.com/sazonovItas/go-tcpws"
)

// Response for the
type Response struct {
	Status     string
	StatusCode int

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

func (resp *Response) write() error {
	type response struct {
		Status     string `json:"status"`
		StatusCode int    `json:"status_code"`

		Header map[string]interface{} `json:"header"`
		Body   string                 `json:"body"`
	}

	resp.Header["Content-Length"] = len(resp.Body)

	wrtResp := response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,

		Header: resp.Header,
		Body:   resp.Body,
	}

	err := json.NewEncoder(resp.Conn).Encode(wrtResp)
	return err
}
