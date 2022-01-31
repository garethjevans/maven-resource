package download

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestAllVersions(t *testing.T) {
	a := Artifact{
		GroupId:       "bing",
		Id:            "bong",
		RepositoryUrl: "https://repo1.maven.org/maven2",
		Downloader: func(url string, user string, password string) (*http.Response, error) {
			assert.Equal(t, url, "https://repo1.maven.org/maven2/bing/bong/maven-metadata.xml")

			r := &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?>
<metadata>
  <groupId>bing</groupId>
  <artifactId>bong</artifactId>
  <versioning>
    <latest>42.3.1</latest>
    <release>42.3.1</release>
    <versions>
      <version>42.0.0</version>
      <version>42.1.0</version>
      <version>42.2.0</version>
      <version>42.3.0</version>
      <version>42.3.1</version>
    </versions>
    <lastUpdated>20211029172338</lastUpdated>
  </versioning>
</metadata>
`)),
			}

			return r, nil
		},
	}
	allVersions, err := AllVersions(a)
	assert.NoError(t, err)
	assert.Equal(t, len(allVersions), 5)
}

func TestDownloadArtifact(t *testing.T) {
	a := Artifact{
		GroupId:       "bing",
		Id:            "bong",
		Version:       "0.1",
		RepositoryUrl: "https://repo1.maven.org/maven2",
		Downloader: func(url string, user string, password string) (*http.Response, error) {
			switch url {
			case "https://repo1.maven.org/maven2/bing/bong/0.1/bong-0.1.jar":
				r := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`some dummy jar content`)),
				}
				return r, nil
			case "https://repo1.maven.org/maven2/bing/bong/0.1/bong-0.1.jar.sha1":
				r := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`sha1 hash`)),
				}
				return r, nil
			case "https://repo1.maven.org/maven2/bing/bong/0.1/bong-0.1.jar.sha256":
				r := &http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(strings.NewReader(`404`)),
				}
				return r, nil
			case "https://repo1.maven.org/maven2/bing/bong/0.1/bong-0.1.jar.sha512":
				r := &http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(strings.NewReader(`404`)),
				}
				return r, nil
			}

			return nil, fmt.Errorf("unknown url %s", url)
		},
	}

	dir, err := ioutil.TempDir(".", "test-output")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	_, err = DownloadArtifact(a, dir)
	assert.NoError(t, err)
}
