package etherscan

import (
	"net/http"
	"sync/atomic"
)

type CircleRoundTripper struct {
	roundTrippers []http.RoundTripper
	counter       uint32
}

func NewCircleRT() *CircleRoundTripper {
	return &CircleRoundTripper{
		roundTrippers: make([]http.RoundTripper, 0),
	}
}

func (crt *CircleRoundTripper) AddRT(rt http.RoundTripper) {
	crt.roundTrippers = append(crt.roundTrippers, rt)
}

func (crt *CircleRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	size := uint32(len(crt.roundTrippers))
	var indexRT uint32
	switch size {
	case 0:
		return nil, &errNotFoundRT{}
	case 1:
		indexRT = 0
	default:
		indexRT = atomic.AddUint32(&crt.counter, 1) % size
	}

	return crt.roundTrippers[indexRT].RoundTrip(req)
}
