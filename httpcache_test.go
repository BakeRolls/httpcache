package httpcache_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BakeRolls/httpcache"
	"github.com/BakeRolls/httpcache/memcache"
)

// Generate a new cached http.Client that saves responses with a status code
// between 200 and 300 in memory,
func ExampleNew() {
	cache := memcache.New()
	client := httpcache.New(cache,
		// httpcache.Verify(httpcache.StatusInTwoHundreds)
		httpcache.Verify(func(res *http.Response) bool {
			return res.StatusCode >= 200 && res.StatusCode < 300
		}),
	).Client()

	res, err := client.Get("https://httpbin.org/ip")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	// Output: 200
}
