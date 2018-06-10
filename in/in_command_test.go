package in_test

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	initializr "github.com/jghiloni/spring-initializr-resource"
	"github.com/jghiloni/spring-initializr-resource/in"
	"github.com/jghiloni/spring-initializr-resource/internal"

	. "github.com/onsi/gomega"
)

func TestInCommand(t *testing.T) {
	spec.Run(t, "In Command", func(t *testing.T, when spec.G, it spec.S) {
		when("Testing the in command", func() {
			var initializrServer *httptest.Server
			var request in.Request
			var command *in.Command

			var destDir string

			it.Before(func() {
				RegisterTestingT(t)

				dataDir, err := filepath.Abs("testdata")
				Expect(err).NotTo(HaveOccurred())

				initializrServer = internal.MockInitializrServer(dataDir)

				serverURL, err := url.Parse(initializrServer.URL)
				Expect(err).NotTo(HaveOccurred())

				request = in.Request{
					Source: initializr.Source{
						URL:               serverURL,
						SkipTLSValidation: true,
					},
					Version: initializr.Version{
						ID: "2.0.2.RELEASE",
					},
					Params: in.Params{
						Type: "maven-build",
					},
				}

				client, err := initializr.NewHTTPClient(request.Source)
				Expect(err).NotTo(HaveOccurred())

				command = &in.Command{
					Client: client,
				}

				tmpPath, err := ioutil.TempDir("", "in_command")
				Expect(err).NotTo(HaveOccurred())

				destDir = filepath.Join(tmpPath, "destination")
			})

			it.After(func() {
				initializrServer.Close()
			})

			it("Should create the destination directory", func() {
				Expect(destDir).NotTo(BeADirectory())

				_, err := command.Run(destDir, request)
				Expect(err).NotTo(HaveOccurred())

				Expect(destDir).To(BeADirectory())
			})

			it("Should download all the files", func() {
				_, err := command.Run(destDir, request)
				Expect(err).NotTo(HaveOccurred())

				Expect(filepath.Join(destDir, "pom.xml")).To(BeARegularFile())
				Expect(filepath.Join(destDir, "version")).To(BeARegularFile())
				Expect(filepath.Join(destDir, "url")).To(BeARegularFile())
				Expect(filepath.Join(destDir, "available-dependencies")).To(BeARegularFile())

				fileList, err := ioutil.ReadDir(destDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(fileList).To(HaveLen(4))
			})

			it("Should generate all the appropriate metadata", func() {
				resp, err := command.Run(destDir, request)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.Metadata[0].Name).To(Equal("file"))
				Expect(resp.Metadata[0].Value).To(Equal("pom.xml"))

				Expect(resp.Metadata[1].Name).To(Equal("version"))
				Expect(resp.Metadata[1].Value).To(Equal("2.0.2.RELEASE"))
			})

			it("Should not corrupt the downloaded files", func() {
				_, err := command.Run(destDir, request)
				Expect(err).NotTo(HaveOccurred())

				pomFile, err := os.Open(filepath.Join(destDir, "pom.xml"))
				Expect(err).NotTo(HaveOccurred())

				versionBytes, err := ioutil.ReadFile(filepath.Join(destDir, "version"))
				Expect(err).NotTo(HaveOccurred())

				urlBytes, err := ioutil.ReadFile(filepath.Join(destDir, "url"))
				Expect(err).NotTo(HaveOccurred())

				depBytes, err := ioutil.ReadFile(filepath.Join(destDir, "available-dependencies"))
				Expect(err).NotTo(HaveOccurred())

				err = xml.NewDecoder(pomFile).Decode(new(interface{}))
				Expect(err).NotTo(HaveOccurred())

				depBody := make([]string, 0, 116)
				err = json.Unmarshal(depBytes, &depBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(depBody).To(HaveLen(116))

				Expect(string(urlBytes)).To(Equal(initializrServer.URL + "/pom.xml?bootVersion=2.0.2.RELEASE&type=maven-build"))
				Expect(string(versionBytes)).To(Equal("2.0.2.RELEASE"))
			})
		})
	}, spec.Report(report.Terminal{}))
}
