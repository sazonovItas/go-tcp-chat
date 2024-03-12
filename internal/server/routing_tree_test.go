package tcpws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePattern(t *testing.T) {
	type req struct {
		url    string
		method string
	}

	tests := map[req]*pattern{
		{
			url:    "/{dist}//article/",
			method: "POST",
		}: {
			str:    "/{dist}//article/",
			method: "POST",
			segments: []segment{
				{
					str:  "dist",
					wild: true,
				},
				{
					str:  "/",
					wild: false,
				},
				{
					str:  "article",
					wild: false,
				},
			},
		},
		{
			url:    "/user/{id}/description",
			method: "GET",
		}: {
			str:    "/user/{id}/description",
			method: "GET",
			segments: []segment{
				{
					str:  "user",
					wild: false,
				},
				{
					str:  "id",
					wild: true,
				},
				{
					str:  "description",
					wild: false,
				},
			},
		},
	}

	for key, value := range tests {
		t.Run(
			fmt.Sprintf("check request with url %s and method %s", key.url, key.method),
			func(t *testing.T) {
				pattern, err := parsePattern(key.method, key.url)
				if err != nil {
					assert.Equal(t, nil, err, "should not be error parsing pattern")
					return
				}
				assert.Equal(t, value, pattern, "should be equal patterns")
			},
		)
	}

	t.Run("check empty pattern", func(t *testing.T) {
		_, err := parsePattern("", "")
		if assert.Error(t, err, "should be error parsing empty pattern") {
			assert.Equal(t, errEmptyPattern, err, "should be empty pattern error")
		}
	})
}

func Test_addPattern_match(t *testing.T) {
	type req struct {
		url    string
		method string
	}

	testReqs := []req{
		{
			url:    "/{dist}/article/",
			method: "POST",
		},
		{
			url:    "/user/{id}/description",
			method: "GET",
		},
		{
			url:    "/{dist}/{articleId}/article/",
			method: "POST",
		},
		{
			url:    "/user/hello",
			method: "GET",
		},
	}

	var patterns []*pattern
	for _, r := range testReqs {
		p, err := parsePattern(r.method, r.url)
		if err != nil {
			assert.Errorf(t, err, "should not be error converting req to pattern")
			continue
		}

		patterns = append(patterns, p)
	}

	root := &routingNode{
		children: map[string]*routingNode{},
	}

	for _, p := range patterns {
		root.addPattern(p, nil)
	}

	t.Run("check GET /user/hello", func(t *testing.T) {
		var (
			n *routingNode
			m []string
		)

		n, m = root.match("GET", "/user/hello")
		assert.Equal(t, patterns[3], n.pattern, "should be equal patterns")
		assert.Equal(t, []string(nil), m, "should not have matches")
	})

	t.Run("check GET /user/31231/description", func(t *testing.T) {
		var (
			n *routingNode
			m []string
		)

		n, m = root.match("GET", "/user/31231/description")
		assert.Equal(t, patterns[1], n.pattern, "should be equals patterns")
		assert.Equal(t, []string{"31231"}, m, "should have mathes")
	})

	t.Run("check POST /123/article", func(t *testing.T) {
		var (
			n *routingNode
			m []string
		)

		n, m = root.match("POST", "/123/article")
		assert.Equal(t, patterns[0], n.pattern, "should be equals patterns")
		assert.Equal(t, []string{"123"}, m, "should have mathes")
	})

	t.Run("check POST /123/article", func(t *testing.T) {
		var (
			n *routingNode
			m []string
		)

		n, m = root.match("POST", "/hi/123/article")
		assert.Equal(t, patterns[2], n.pattern, "should be equals patterns")
		assert.Equal(t, []string{"hi", "123"}, m, "should have mathes")
	})

	t.Run("check set node panic", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotEqual(t, nil, r, "should recover from panic %s", r)
		}()

		root.addPattern(patterns[0], nil)
	})
}
