package download

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
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

func Download(groupId, artifactId, version, dest, repo, filename, extension, user, pwd string) (string, string, error) {
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

	url, err := ArtifactUrl(a)
	if err != nil {
		return "", "", err
	}
	//
	//resp, err := a.Downloader(url, user, pwd)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()
	//
	//if filename == "" {
	//	filename = FileName(a)
	//}
	//
	//filepath := dest + "/" + filename
	//
	//out, err := os.Create(filepath)
	//if err != nil {
	//	return "", err
	//}
	//defer out.Close()
	//
	//_, err = io.Copy(out, resp.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//return filepath, nil
	return a.Version, url, nil
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

func ArtifactUrl(a Artifact) (string, error) {
	if a.RepositoryUrl == "" {
		a.RepositoryUrl = "https://repo1.maven.org/maven2/"
	}

	//if a.IsSnapshot {
	//	var err error
	//	a.SnapshotVersion, err = LatestSnapshotVersion(a)
	//	if err != nil {
	//		return "", err
	//	}
	//}

	// FIXME should ensure that repo url has a trailing slash
	return a.RepositoryUrl + "/" + artifactPath(a), nil
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
