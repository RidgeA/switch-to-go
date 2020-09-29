package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
)

type (
	Cache interface {
		Set(string, interface{})
		Get(string) (interface{}, bool)
		Del(string)
	}

	ResponseCacheEntry struct {
		Body    []byte
		Status  int
		Headers http.Header
	}

	responseCacher struct {
		http.ResponseWriter
		body   *bytes.Buffer
		status int
	}
)

func (bc *responseCacher) Write(b []byte) (int, error) {
	if bc.body == nil {
		bc.body = bytes.NewBuffer([]byte{})
	}
	bc.body.Write(b)
	return bc.ResponseWriter.Write(b)
}

func (bc *responseCacher) WriteHeader(status int) {
	bc.status = status
	bc.ResponseWriter.WriteHeader(status)
}

func (bc responseCacher) Cache() ResponseCacheEntry {
	rce := ResponseCacheEntry{Status: bc.status, Headers: bc.Header().Clone()}

	if bc.body != nil {
		rce.Body = bc.body.Bytes()
	}
	return rce
}

func NewCacher(c Cache) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			key := reqKey(req)

			if data, exists := c.Get(key); exists {
				cachedResponse, ok := data.(ResponseCacheEntry)
				if !ok {
					c.Del(key)
				} else {
					returnCachedResponse(cachedResponse, res)
					return
				}
			}

			bc := wrapWithResponseCacher(res)

			next(bc, req)

			c.Set(key, bc.Cache())
		}
	}
}

func reqKey(req *http.Request) string {
	return fmt.Sprintf("%s%s", req.Method, req.URL.String())
}

func returnCachedResponse(cacheEntity ResponseCacheEntry, res http.ResponseWriter) {
	res.WriteHeader(cacheEntity.Status)
	for key, values := range res.Header() {
		for _, value := range values {
			res.Header().Add(key, value)
		}
	}
	if cacheEntity.Body != nil {
		_, _ = res.Write(cacheEntity.Body)
	}
}

func wrapWithResponseCacher(rw http.ResponseWriter) *responseCacher {
	return &responseCacher{rw, nil, http.StatusOK}
}
