package http

import (
	"net/http"
)

type Requester func(req *http.Request) (*http.Response, error)

type ClientInterceptor func(req *http.Request, requester Requester) (*http.Response, error)

type Client struct {
	interceptor ClientInterceptor
	client      *http.Client
}

type config struct {
	http.Client
	interceptors []ClientInterceptor
}

func WithClientInterceptors(interceptors ...ClientInterceptor) Option {
	return func(c *config) {
		c.interceptors = append(c.interceptors, interceptors...)
	}
}

type Option func(*config)

func NewClient(opts ...Option) Client {
	c := config{}

	for _, opt := range opts {
		opt(&c)
	}
	return Client{
		interceptor: chainedClientInterceptor(c.interceptors),
		client:      &http.Client{},
	}
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.interceptor(req, func(req *http.Request) (*http.Response, error) {
		return c.client.Do(req)
	})
}

func genChainClientHandler(interceptors []ClientInterceptor, cur int, requester Requester) Requester {
	if cur >= len(interceptors)-1 {
		return requester
	}
	return func(req *http.Request) (*http.Response, error) {
		return interceptors[cur+1](req, genChainClientHandler(interceptors, cur+1, requester))
	}
}

func chainedClientInterceptor(chainedInterceptors []ClientInterceptor) ClientInterceptor {
	iNums := len(chainedInterceptors)
	if iNums == 0 {
		return func(req *http.Request, requester Requester) (*http.Response, error) {
			return requester(req)
		}
	} else if iNums == 1 {
		return chainedInterceptors[0]
	} else {
		return func(req *http.Request, requester Requester) (*http.Response, error) {
			return chainedInterceptors[0](req, genChainClientHandler(chainedInterceptors, 0, requester))
		}
	}
	return nil
}
