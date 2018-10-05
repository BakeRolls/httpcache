package httpcache

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
)

// Cache is a key value store.
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, dump []byte) error
}

// Transport implements the http.RoundTripper interface.
type Transport struct {
	cache  Cache
	verify func(*http.Response) bool
}

// StatusInTwoHundreds returns true if the responses' status code is between
// 200 and 300.
func StatusInTwoHundreds(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

// Verify checks if a http.Response is cachable, i.e. the status code is OK or
// the body contains relevant information.
func Verify(fn func(*http.Response) bool) func(*Transport) {
	return func(t *Transport) {
		t.verify = fn
	}
}

// New returns a new cachable Transport.
func New(c Cache, options ...func(*Transport)) *Transport {
	t := &Transport{c, func(*http.Response) bool { return true }}
	for _, option := range options {
		option(t)
	}
	return t
}

// RoundTrip executes a single HTTP transaction, returning a Response for the
// provided Request.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()
	if dump, err := t.cache.Get(key); err == nil {
		buf := bufio.NewReader(bytes.NewReader(dump))
		return http.ReadResponse(buf, req)
	}

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if !t.verify(res) {
		return res, nil
	}
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	if err := t.cache.Set(key, dump); err != nil {
		return nil, err
	}
	return res, nil
}

// Client returns a new cached http.Client.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}
