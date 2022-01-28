package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/garethjevans/maven-resource/download"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type InCmd struct {
	Command *cobra.Command
}

func NewInCmd() InCmd {
	in := InCmd{}
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

	version, uri, err := download.Download(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, jsonIn.Version.Ref,
		".", jsonIn.Source.Repository, "download.jar", jsonIn.Source.Type, jsonIn.Source.Username, jsonIn.Source.Password)
	if err != nil {
		panic(err)
	}

	out := InResponse{
		Version: Version{Ref: version},
		Metadata: []Metadata{
			{Name: "version", Value: version},
			{Name: "uri", Value: uri},
		},
	}

	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
