package check_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/jghiloni/spring-initializr-resource"

	"github.com/jghiloni/spring-initializr-resource/check"
	"github.com/jghiloni/spring-initializr-resource/internal"

	. "github.com/onsi/gomega"
)

func TestCheckCommand(t *testing.T) {
	spec.Run(t, "Check Command", func(t *testing.T, when spec.G, it spec.S) {
		when("Testing the check command", func() {
			var request check.Request

			var initializrServer *httptest.Server

			it.Before(func() {
				RegisterTestingT(t)

				dataDir, err := filepath.Abs("testdata")
				Expect(err).NotTo(HaveOccurred())

				initializrServer = internal.MockInitializrServer(dataDir)
			})

			it.After(func() {
				initializrServer.Close()
			})

			when("I test basic check functionality", func() {
				var fakeClient *http.Client

				it.Before(func() {
					var err error
					serverURL, err := url.Parse(initializrServer.URL)
					Expect(err).NotTo(HaveOccurred())

					baseSource := initializr.Source{
						URL:               serverURL,
						SkipTLSValidation: true,
					}

					fakeClient, err = initializr.NewHTTPClient(baseSource)
					Expect(err).NotTo(HaveOccurred())
				})

				when("I get versions for the first time", func() {
					it("returns all the release versions", func() {
						bytes, err := ioutil.ReadFile("testdata/first_request.json")
						Expect(err).NotTo(HaveOccurred())

						err = json.Unmarshal(bytes, &request)
						Expect(err).NotTo(HaveOccurred())

						cmd := &check.Command{
							Client: fakeClient,
						}

						resp, err := cmd.Run(request)
						Expect(err).NotTo(HaveOccurred())
						Expect(resp).To(HaveLen(2))
						Expect(resp[0].ID).To(Equal("2.0.2.RELEASE"))
					})

					it("returns all the versions", func() {
						bytes, err := ioutil.ReadFile("testdata/first_request_with_snapshots.json")
						Expect(err).NotTo(HaveOccurred())

						err = json.Unmarshal(bytes, &request)
						Expect(err).NotTo(HaveOccurred())

						cmd := &check.Command{
							Client: fakeClient,
						}

						resp, err := cmd.Run(request)
						Expect(err).NotTo(HaveOccurred())
						Expect(resp).To(HaveLen(5))
						Expect(resp[0].ID).To(Equal("2.1.0.BUILD-SNAPSHOT"))
					})
				})

				when("I have checked recently", func() {
					it("returns only the latest version", func() {
						bytes, err := ioutil.ReadFile("testdata/subsequent_request.json")
						Expect(err).NotTo(HaveOccurred())

						err = json.Unmarshal(bytes, &request)
						Expect(err).NotTo(HaveOccurred())

						cmd := &check.Command{
							Client: fakeClient,
						}

						resp, err := cmd.Run(request)
						Expect(err).NotTo(HaveOccurred())
						Expect(resp).To(HaveLen(1))
						Expect(resp[0].ID).To(Equal("2.0.2.RELEASE"))
					})

					it("returns all later versions", func() {
						bytes, err := ioutil.ReadFile("testdata/subsequent_request_with_snapshots.json")
						Expect(err).NotTo(HaveOccurred())

						err = json.Unmarshal(bytes, &request)
						Expect(err).NotTo(HaveOccurred())

						cmd := &check.Command{
							Client: fakeClient,
						}

						resp, err := cmd.Run(request)
						Expect(err).NotTo(HaveOccurred())
						Expect(resp).To(HaveLen(3))
						Expect(resp[0].ID).To(Equal("2.1.0.BUILD-SNAPSHOT"))
					})
				})

				when("I have pinned to a specific major minor version", func() {
					it("returns no versions", func() {
						bytes, err := ioutil.ReadFile("testdata/subsequent_request_with_pin.json")
						Expect(err).NotTo(HaveOccurred())

						err = json.Unmarshal(bytes, &request)
						Expect(err).NotTo(HaveOccurred())

						cmd := &check.Command{
							Client: fakeClient,
						}

						resp, err := cmd.Run(request)
						Expect(err).NotTo(HaveOccurred())
						Expect(resp).To(HaveLen(0))
					})

					it("returns a version with snapshots enabled", func() {
						bytes, err := ioutil.ReadFile("testdata/subsequent_request_with_pin_and_snapshots.json")
						Expect(err).NotTo(HaveOccurred())

						err = json.Unmarshal(bytes, &request)
						Expect(err).NotTo(HaveOccurred())

						cmd := &check.Command{
							Client: fakeClient,
						}

						resp, err := cmd.Run(request)
						Expect(err).NotTo(HaveOccurred())
						Expect(resp).To(HaveLen(1))
						Expect(resp[0].ID).To(Equal("1.5.14.BUILD-SNAPSHOT"))
					})
				})
			})
		})
	}, spec.Report(report.Terminal{}))
}
