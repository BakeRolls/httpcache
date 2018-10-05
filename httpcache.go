package httpcache

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, dump []byte) error
}

type Transport struct{ Cache }

func New(c Cache) *Transport {
	return &Transport{c}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()
	if dump, err := t.Get(key); err == nil {
		buf := bufio.NewReader(bytes.NewReader(dump))
		return http.ReadResponse(buf, req)
	}

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	if err := t.Set(key, dump); err != nil {
		return nil, err
	}
	return res, nil
}

func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}
