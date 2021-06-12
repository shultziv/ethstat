package config

import (
	"net/url"
)

type EtherScan struct {
	Proxies []string `env:"PROXIES" envSeparator:";"`
}

func (e EtherScan) GetProxyURLs() (proxyUrls []*url.URL, err error) {
	proxyUrls = make([]*url.URL, len(e.Proxies))
	for i, proxyUrl := range e.Proxies {
		proxyUrls[i], err = url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
	}

	return proxyUrls, nil
}
