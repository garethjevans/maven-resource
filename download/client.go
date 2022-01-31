package download

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Artifact struct {
	GroupId       string
	Id            string
	Version       string
	Extension     string
	Classifier    string
	RepositoryUrl string
	RepoUser      string
	RepoPassword  string
	Downloader    func(string, string, string) (*http.Response, error)
}

type DownloadedArtifact struct {
	Filename string
	Url      string
	Version  string
	Sha1     string
	Sha256   string
}

type metadata struct {
	Versions []string `xml:"versioning>versions>version"`
}

func GetVersions(groupId, artifactId, repository, username, password string) ([]string, error) {
	a := Artifact{
		GroupId:       groupId,
		Id:            artifactId,
		RepositoryUrl: repository,
		Downloader:    httpGetCustom,
		RepoUser:      username,
		RepoPassword:  password,
	}

	v, err := AllVersions(a)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func Download(groupId, artifactId, version, dest, repo, extension, user, pwd string) (*DownloadedArtifact, error) {
	a := Artifact{
		GroupId:       groupId,
		Id:            artifactId,
		Extension:     extension,
		Version:       version,
		RepositoryUrl: repo,
		Downloader:    httpGetCustom,
		RepoUser:      user,
		RepoPassword:  pwd,
	}

	return DownloadArtifact(a, dest)
}

func DownloadArtifact(a Artifact, dest string) (*DownloadedArtifact, error) {
	url := ArtifactUrl(a)

	resp, err := a.Downloader(url, a.RepoUser, a.RepoPassword)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	filename := FileName(a)
	filepath := dest + "/" + filename

	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	sha1 := Sha1(a)
	sha256 := Sha256(a)

	return &DownloadedArtifact{Version: a.Version, Url: url, Filename: filename, Sha1: sha1, Sha256: sha256}, nil
}

func httpGetCustom(url, user, pwd string) (*http.Response, error) {
	if user != "" && pwd != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(user, pwd)
		return client.Do(req)
	}

	return http.Get(url)
}

func FileName(a Artifact) string {
	ext := "jar"
	if a.Extension != "" {
		ext = a.Extension
	}

	v := a.Version

	if a.Classifier != "" {
		return fmt.Sprintf("%s-%s-%s.%s", a.Id, v, a.Classifier, ext)
	} else {
		return fmt.Sprintf("%s-%s.%s", a.Id, v, ext)
	}
}

func ArtifactUrl(a Artifact) string {
	if a.RepositoryUrl == "" {
		a.RepositoryUrl = "https://repo1.maven.org/maven2/"
	}

	// FIXME should ensure that repo url has a trailing slash
	return a.RepositoryUrl + "/" + artifactPath(a)
}

func Sha1(a Artifact) string {
	url := ArtifactUrl(a) + ".sha1"
	r, err := a.Downloader(url, a.RepoUser, a.RepoPassword)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func Sha256(a Artifact) string {
	url := ArtifactUrl(a) + ".sha256"
	r, err := a.Downloader(url, a.RepoUser, a.RepoPassword)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func AllVersions(a Artifact) ([]string, error) {
	// FIXME should ensure that repo url has a trailing slash
	metadataUrl := a.RepositoryUrl + "/" + groupPath(a) + "/maven-metadata.xml"
	resp, err := a.Downloader(metadataUrl, a.RepoUser, a.RepoPassword)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unable to fetch maven metadata from %s Http statusCode: %d", metadataUrl, resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := metadata{}
	err = xml.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	return m.Versions, nil
}

func artifactPath(a Artifact) string {
	return groupPath(a) + "/" + FileName(a)
}

func groupPath(a Artifact) string {
	parts := append(strings.Split(a.GroupId, "."), a.Id)
	if a.Version != "" {
		parts = append(parts, a.Version)
	}
	return strings.Join(parts, "/")
}
