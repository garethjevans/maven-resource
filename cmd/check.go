package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/garethjevans/maven-resource/download"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type CheckCmd struct {
	Command *cobra.Command
}

func NewCheckCmd() CheckCmd {
	check := CheckCmd{}
	check.Command = &cobra.Command{
		Use:   "check",
		Short: "checks a resource",
		Long:  `checks a resource`,
		Run:   check.Run,
	}

	return check
}

func (i *CheckCmd) Run(cmd *cobra.Command, args []string) {
	var jsonIn In

	err := json.NewDecoder(os.Stdin).Decode(&jsonIn)
	if err != nil {
		log.Fatal(err)
	}

	version, _, err := download.Download(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, "", ".", jsonIn.Source.Repository, "download.jar", jsonIn.Source.Type, "", "")
	if err != nil {
		panic(err)
	}

	refs := []Version{{Ref: version}}
	b, err := json.Marshal(refs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
