package initializr

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

// UnmarshalJSON unmarshals and verifies the source json block
func (s *Source) UnmarshalJSON(j []byte) error {
	intermediate := make(map[string]interface{})
	var err error

	if err = json.Unmarshal(j, &intermediate); err != nil {
		return err
	}

	// set defaults
	if _, ok := intermediate["url"]; !ok {
		intermediate["url"] = "https://start.spring.io"
	}

	for key, val := range intermediate {
		switch key {
		case "url":
			if s.URL, err = makeURL(val); err != nil {
				return err
			}
		case "skip_tls_validation":
			if s.SkipTLSValidation, err = makeBool(val); err != nil {
				return err
			}
		case "product_version":
			if re, ok := val.(string); ok {
				if s.ProductVersion, err = regexp.Compile(re); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("product_version needs to be a string, got a %T", val)
			}
		case "include_snapshots":
			if s.IncludeSnapshots, err = makeBool(val); err != nil {
				return err
			}
		case "ca_certs":
			if _, ok := val.([]string); ok {
				for _, certPEM := range val.([]string) {
					var cert *x509.Certificate
					if cert, err = makeCertificate(certPEM); err != nil {
						return err
					}

					if s.CACerts == nil {
						s.CACerts = make([]*x509.Certificate, 0, 3)
					}

					s.CACerts = append(s.CACerts, cert)
				}
			}

			return fmt.Errorf("ca_certs must be a list of strings, got a %T", val)
		case "http_proxy":
			if _, ok := val.(string); ok {
				s.HTTPProxy = val.(string)
			}

			return fmt.Errorf("http_proxy must be a string, got a %T", val)
		case "https_proxy":
			if _, ok := val.(string); ok {
				s.HTTPSProxy = val.(string)
			}

			return fmt.Errorf("https_proxy must be a string, got a %T", val)
		case "no_proxy":
			if _, ok := val.(string); ok {
				s.NoProxy = val.(string)
			}

			return fmt.Errorf("no_proxy must be a string, got a %T", val)
		default:
			return fmt.Errorf("Field is %s ... should it be under source?", key)
		}
	}

	return nil
}

func makeURL(val interface{}) (*url.URL, error) {
	var newURL *url.URL
	var err error
	if rawurl, ok := val.(string); ok {
		if newURL, err = url.Parse(rawurl); err != nil {
			return nil, err
		}

		return newURL, nil
	}

	return nil, fmt.Errorf("Expected URL to be a string, got a %T instead", val)
}

func makeBool(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}

	if newBool, ok := val.(bool); ok {
		return newBool, nil
	}

	if maybeBoolStr, ok := val.(string); ok {
		return strconv.ParseBool(maybeBoolStr)
	}

	return false, fmt.Errorf("Expected bool or boolean string, got a %T instead", val)
}

func makeCertificate(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errors.New("Could not decode certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}
