package initializr

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// NewHTTPClient will create an HTTP client configured with the SSL options
func NewHTTPClient(source Source) (*http.Client, error) {
	certs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if len(source.CACerts) > 0 {
		for i := range source.CACerts {
			certs.AddCert(source.CACerts[i])
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            certs,
			},
			Proxy: func(req *http.Request) (*url.URL, error) {
				if strings.TrimSpace(source.HTTPProxy) != "" {
					os.Setenv("HTTP_PROXY", source.HTTPProxy)
				}

				if strings.TrimSpace(source.HTTPSProxy) != "" {
					os.Setenv("HTTPS_PROXY", source.HTTPSProxy)
				}

				if strings.TrimSpace(source.NoProxy) != "" {
					os.Setenv("NO_PROXY", source.NoProxy)
				}

				return http.ProxyFromEnvironment(req)
			},
		},
	}, nil
}
