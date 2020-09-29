package middlewares

import (
	"log"
	"net/http"
	"time"
)

type logged struct {
	http.ResponseWriter
	statusCode int
}

func (l *logged) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func NewLogger(l *log.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {

		return func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()
			resWithLog := wrapWithLogged(res)

			next(resWithLog, req)

			l.Printf(" %s [%s] %s - %d - %dms",
				req.RemoteAddr,
				start.Format(time.RFC3339),
				req.URL.String(),
				resWithLog.statusCode,
				time.Now().Sub(start).Milliseconds(),
			)
		}
	}
}

func wrapWithLogged(rw http.ResponseWriter) *logged {
	return &logged{rw, http.StatusOK}
}
