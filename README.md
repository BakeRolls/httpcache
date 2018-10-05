# httpcache

[![GoDoc](https://godoc.org/github.com/BakeRolls/httpcache?status.svg)](https://godoc.org/github.com/BakeRolls/httpcache)
[![Go Report Card](https://goreportcard.com/badge/github.com/BakeRolls/httpcache)](https://goreportcard.com/report/github.com/BakeRolls/httpcache)

As simple as it gets HTTP cache. It is heavily inspired by [gregjones/httpcache](https://github.com/gregjones/httpcache), but it ignores the response headers. If you need a cache that respects these, please take a look at [gregjones/httpcache](https://github.com/gregjones/httpcache) and [lox/httpcache](https://github.com/lox/httpcache).

## Usage

```go
cache := memcache.New()
httpClient := httpcache.New(cache).Client()
```

Per default every response is cached, even ones with a status code outside the range [200-300[. You can specify if a response should get saved by providing a function using the `Verify` option.

```go
cache := memcache.New()
client := httpcache.New(cache,
	httpcache.Verify(func(res *http.Response) bool {
		return res.StatusCode >= 200 && res.StatusCode < 300
	}),
).Client()
```

If the mentioned range is verification enough, a function `httpcache.StatusInTwoHundreds` is provided.

## Cache

Currently there are two cache methods: in memory and on disk using [diskv](https://github.com/peterbourgon/diskv). Caches have to implement the [`Cache` interface](https://godoc.org/github.com/BakeRolls/httpcache#Cache), which is basically a key-value-store with two functions `Get` and `Set`.
