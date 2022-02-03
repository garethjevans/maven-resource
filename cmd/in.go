package cmd

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/garethjevans/maven-resource/download"
	"github.com/spf13/cobra"
)

type InCmd struct {
	Command    *cobra.Command
	Downloader download.Downloader
}

func NewInCmd() InCmd {
	in := InCmd{Downloader: &download.DefaultDownloader{}}
	in.Command = &cobra.Command{
		Use:   "in",
		Short: "concourse in command",
		Long:  `concourse in command`,
		Run:   in.Run,
	}

	return in
}

func (i *InCmd) Run(cmd *cobra.Command, args []string) {
	var jsonIn In

	err := json.NewDecoder(os.Stdin).Decode(&jsonIn)
	if err != nil {
		log.Fatal(err)
	}

	outputDir := args[0]

	artifact, err := i.Downloader.Download(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, jsonIn.Version.Ref, outputDir, jsonIn.Source.Repository, jsonIn.Source.Type)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// lets validate sha1, this should always exist
	downloadedFilePath := path.Join(outputDir, artifact.Filename)
	downloadedFileContents, err := ioutil.ReadFile(downloadedFilePath)
	if err != nil {
		panic(err)
	}

	sha1 := sha1.Sum(downloadedFileContents)

	if fmt.Sprintf("%x", sha1) != artifact.Sha1 {
		log.Fatalf("calculated sha1 does not match downloaded sha1: %x != %s\n", sha1, artifact.Sha1)
	} else {
		Log("sha1 %s is valid\n", artifact.Sha1)
	}

	// if sha256 does exist, calculate it
	if artifact.Sha256 == "" {
		Log("sha256 does not exist, calculating it from downloaded file\n")
		sha256 := sha256.Sum256(downloadedFileContents)
		artifact.Sha256 = fmt.Sprintf("%x", sha256)
	}

	versionWithoutRelease := strings.Replace(artifact.Version, ".RELEASE", "", -1)

	out := InResponse{
		Version: Version{Ref: artifact.Version},
		Metadata: []Metadata{
			{Name: "version", Value: versionWithoutRelease},
			{Name: "uri", Value: artifact.Url},
			{Name: "filename", Value: artifact.Filename},
			{Name: "cpe", Value: versionWithoutRelease},
			{Name: "purl", Value: versionWithoutRelease},
			{Name: "sha1", Value: artifact.Sha1},
			{Name: "sha256", Value: artifact.Sha256},
		},
	}

	for _, m := range out.Metadata {
		file := path.Join(outputDir, m.Name)
		Log("creating %s\n", file)
		err = ioutil.WriteFile(file, []byte(m.Value), 0644)
		if err != nil {
			panic(err)
		}
	}

	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
