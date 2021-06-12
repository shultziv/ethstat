package etherscan

import (
	"net/http"
	"sync/atomic"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CircleHttpClient запускает внутренние клиенты по циклу
type CircleHttpClient struct {
	*http.Client
	clients []httpClient
	counter uint32
}

func NewCircleHC() *CircleHttpClient {
	return &CircleHttpClient{
		clients: make([]httpClient, 0),
	}
}

func (crt *CircleHttpClient) AddHC(hc httpClient) {
	crt.clients = append(crt.clients, hc)
}

func (crt *CircleHttpClient) Do(req *http.Request) (res *http.Response, e error) {
	size := uint32(len(crt.clients))
	var indexRT uint32
	switch size {
	case 0:
		return nil, &errNotFoundRT{}
	case 1:
		indexRT = 0
	default:
		indexRT = atomic.AddUint32(&crt.counter, 1) % size
	}

	return crt.clients[indexRT].Do(req)
}
