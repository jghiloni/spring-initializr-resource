package initializr

import (
	"crypto/x509"
	"net/url"
	"regexp"
)

// Source is the data that is defined in the Concourse resource block
type Source struct {
	URL               *url.URL            `json:"url,omitempty"`
	SkipTLSValidation bool                `json:"skip_tls_validation,omitempty"`
	CACerts           []*x509.Certificate `json:"ca_certs,omitempty"`
	ProductVersion    *regexp.Regexp      `json:"product_version,omitempty"`
	IncludeSnapshots  bool                `json:"include_snapshots,omitempty"`
	HTTPProxy         string              `json:"http_proxy,omitempty"`
	HTTPSProxy        string              `json:"https_proxy,omitempty"`
	NoProxy           string              `json:"no_proxy,omitempty"`
}

// Version is the data structure that is output by the check and in scripts
type Version struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// MetadataPair is the datastructure that gets output with a Version from the in script
type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
