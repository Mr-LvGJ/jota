package http

import (
	"net/http"
)

type ServerInterceptor func(w http.ResponseWriter, r *http.Request, handler http.Handler)

func WrapHandler(serverInts []ServerInterceptor, handler http.Handler) http.Handler {
	serverInt := chainedServerInterceptor(serverInts)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverInt(w, r, handler)
	})
}

func chainedServerInterceptor(chainedInts []ServerInterceptor) (chainedInt ServerInterceptor) {
	l := len(chainedInts)
	if l == 0 {
		chainedInt = func(w http.ResponseWriter, r *http.Request, handler http.Handler) {
			handler.ServeHTTP(w, r)
		}
	} else if l == 1 {
		chainedInt = chainedInts[0]
	} else {
		chainedInt = func(w http.ResponseWriter, r *http.Request, handler http.Handler) {
			chainedInts[0](w, r, getChainServerHandler(chainedInts, 0, handler))
		}
	}
	return
}

func getChainServerHandler(interceptors []ServerInterceptor, cur int, handler http.Handler) http.Handler {
	if cur >= len(interceptors)-1 {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := w.(http.ResponseWriter)
		interceptors[cur+1](rw, r, getChainServerHandler(interceptors, cur+1, handler))
	})
}
