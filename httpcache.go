package httpcache

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

// Verifier determine if a given request-response-pair should get cached.
type Verifier func(*http.Request, *http.Response) bool

// Cache is a key value store.
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, dump []byte) error
}

// Transport implements the http.RoundTripper interface.
type Transport struct {
	cache     Cache
	verifiers []Verifier
	transport http.RoundTripper
}

// StatusInTwoHundreds returns true if the responses' status code is between
// 200 and 300.
func StatusInTwoHundreds(req *http.Request, res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

// RequestMethod returns true if a given request method is used.
func RequestMethod(method string) Verifier {
	return func(req *http.Request, res *http.Response) bool {
		return req.Method == method
	}
}

// Verify checks if a http.Response is cachable, i.e. the status code is OK or
// the body contains relevant information.
func Verify(v Verifier) func(*Transport) {
	return func(t *Transport) {
		t.verifiers = append(t.verifiers, v)
	}
}

// WithTransport replaces the http.DefaultTransport RoundTripper.
func WithTransport(transport http.RoundTripper) func(*Transport) {
	return func(t *Transport) {
		t.transport = transport
	}
}

// New returns a new cachable Transport.
func New(c Cache, options ...func(*Transport)) *Transport {
	t := &Transport{cache: c, transport: http.DefaultTransport}
	for _, option := range options {
		option(t)
	}
	return t
}

// RoundTrip executes a single HTTP transaction, returning a Response for the
// provided Request.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()

	if req.Method == http.MethodPost && req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, errors.Wrap(err, "could not read request body")
		}
		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		hash := sha256.Sum256(body)
		key += " " + hex.EncodeToString(hash[:])
	}

	if dump, err := t.cache.Get(key); err == nil {
		buf := bufio.NewReader(bytes.NewReader(dump))
		return http.ReadResponse(buf, req)
	}

	res, err := t.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	for _, verify := range t.verifiers {
		if !verify(req, res) {
			return res, nil
		}
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
