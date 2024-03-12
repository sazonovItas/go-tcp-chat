package tcpws

import (
	"errors"
	"strings"
)

var errEmptyPattern = errors.New("empty pattern")

// pattern is something that mathes a path for request
type pattern struct {
	str    string
	method string

	segments []segment
}

// segment is part of the pattern or url
// /x => {str: x, wild: false}
// /{x} => {str: x, wild true}
type segment struct {
	str  string
	wild bool
}

// parsePattern is creating pattern for given method and url
func parsePattern(method string, url string) (*pattern, error) {
	if len(method) == 0 || len(url) == 0 {
		return nil, errEmptyPattern
	}

	fixUrl := strings.TrimFunc(url, func(c rune) bool {
		return c == ' ' || c == '/'
	})

	segmentsStr := strings.Split(fixUrl, "/")

	var segments []segment
	for _, seg := range segmentsStr {
		seg = strings.TrimFunc(seg, func(c rune) bool {
			return c == ' '
		})

		switch {
		case len(seg) == 0:
			segments = append(segments, segment{
				str:  "/",
				wild: false,
			})
		case seg[0] == '{' && seg[len(seg)-1] == '}':
			str := strings.TrimFunc(seg, func(c rune) bool {
				return c == '{' || c == '}'
			})
			if len(str) == 0 {
				str = "$"
			}
			segments = append(segments, segment{
				str:  str,
				wild: true,
			})
		default:
			segments = append(segments, segment{
				str:  seg,
				wild: false,
			})
		}
	}

	return &pattern{
		str:      url,
		method:   method,
		segments: segments,
	}, nil
}

type routingNode struct {
	// A leaf node holds a single pattern and the HandlerFunc it was registered
	pattern *pattern
	handler HandlerFunc

	// special children keys:
	//    "/" - empty string pattern
	//    "*" - for wild segment
	children map[string]*routingNode
}

// addPattern adds pattern to a tree with handler h with Handler at root
func (root *routingNode) addPattern(p *pattern, h HandlerFunc) {
	n := root.addChild(p.method)

	n.addSegments(p.segments, p, h)
}

func (n *routingNode) addSegments(segs []segment, p *pattern, h HandlerFunc) {
	if len(segs) == 0 {
		n.set(p, h)
		return
	}

	seg := segs[0]
	if seg.wild {
		n.addChild("*").addSegments(segs[1:], p, h)
	} else {
		n.addChild(seg.str).addSegments(segs[1:], p, h)
	}
}

// addChild add children to the node
func (n *routingNode) addChild(key string) *routingNode {
	if c := n.findChild(key); c != nil {
		return c
	}

	c := &routingNode{children: map[string]*routingNode{}}
	n.children[key] = c
	return c
}

// findChild finds child or a routeNode in children map
func (n *routingNode) findChild(key string) *routingNode {
	if r, ok := n.children[key]; ok {
		return r
	}

	return nil
}

// set sets a pattern and handler for the node, it must be leaf node
// if to add pattern that was added then throw panic
func (n *routingNode) set(p *pattern, h HandlerFunc) {
	if n.pattern != nil || n.handler != nil {
		panic("error non-nil leaf pattern or handler")
	}

	n.pattern, n.handler = p, h
}

// match matches params and handler function for a method and url
func (root *routingNode) match(method, url string) (*routingNode, []string) {
	if root == nil {
		return nil, nil
	}

	p, err := parsePattern(method, url)
	if err != nil {
		return nil, nil
	}

	return root.findChild(p.method).matchPath(p.segments, nil)
}

// matchPath mathes params and handler for a segments of pattern
func (n *routingNode) matchPath(segments []segment, matches []string) (*routingNode, []string) {
	if n == nil {
		return nil, nil
	}

	if len(segments) == 0 {
		if n.pattern == nil {
			return nil, nil
		}
		return n, matches
	}

	seg := segments[0]
	if l, m := n.findChild(seg.str).matchPath(segments[1:], matches); l != nil {
		return l, m
	}

	// check for wild key
	if c := n.findChild("*"); c != nil {
		matches = append(matches, seg.str)
		if l, m := c.matchPath(segments[1:], matches); l != nil {
			return l, m
		}

	}

	return nil, nil
}
