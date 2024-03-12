package tcpws

import (
	"context"
	"encoding/json"
)

// A Request represent a custom request received by a server
type Request struct {
	// Method specifies a custom method
	Method string

	// Url and Proto specifies for what this request will be used
	Url   string
	Proto string

	// pattern specifies chosen pattern for request url
	pattern *pattern

	// matches specifies params in url
	matches []string

	// Header specifies some values for the server e.g. session key
	Header map[string]interface{}

	// Body specifies data for the server
	Body string

	// ctx of the request specifies context of the request
	ctx context.Context
}

// Create new request from bytes
func newRequest(msg []byte) (*Request, error) {
	type request struct {
		Method string `json:"method"`
		Url    string `json:"url"`
		Proto  string `json:"proto"`

		Header map[string]interface{} `json:"header"`
		Body   string                 `json:"body"`
	}

	var req request
	err := json.Unmarshal([]byte(string(msg)), &req)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method: req.Method,
		Url:    req.Url,
		Proto:  req.Proto,
		Header: req.Header,
		Body:   req.Body,
	}, nil
}

func (r *Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}

	return r.ctx
}

func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	r2 := new(Request)
	*r2 = *r
	r2.ctx = ctx
	return r2
}

func (r *Request) Params() map[string]string {
	params := make(map[string]string)

	i := 0
	for _, p := range r.pattern.segments {
		if p.wild {
			params[p.str] = r.matches[i]
			i++
		}
	}

	return params
}

func (r *Request) ParamByName(name string) string {
	i := 0
	for _, p := range r.pattern.segments {
		if p.wild {
			if p.str == name {
				return r.matches[i]
			}
			i++
		}
	}

	return ""
}
