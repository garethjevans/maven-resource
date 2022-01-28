package cmd

type Source struct {
	GroupId    string `json:"groupId"`
	ArtifactId string `json:"artifactId"`
	Repository string `json:"repository"`
	Type       string `json:"type"`
	Classifier string `json:"classifier"`
}

type Version struct {
	Ref string `json:"ref"`
}

type In struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
