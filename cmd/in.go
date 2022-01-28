package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/garethjevans/maven-resource/download"

	"github.com/spf13/cobra"
)

type InCmd struct {
	Command    *cobra.Command
	GroupId    string
	ArtifactId string
	Repository string
	Type       string
	Classifier string
}

type Out struct {
	Version string `json:"version"`
	//Sha256  string `json:"sha256"`
	//Cpe     string `json:"cpe"`
	Uri string `json:"uri"`
	// Cpe Pattern
	// PURL
	// PURL Pattern
}

func NewInCmd() InCmd {
	in := InCmd{}
	in.Command = &cobra.Command{
		Use:   "in",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: in.Run,
	}

	//  name, shorthand string, value string, usage string) {
	in.Command.Flags().StringVarP(&in.GroupId, "groupId", "g", "", "The groupId of the artifact to query")
	in.Command.Flags().StringVarP(&in.ArtifactId, "artifactId", "a", "", "The artifactId of the artifact to query")
	in.Command.Flags().StringVarP(&in.Type, "type", "t", "jar", "The type of the artifact to query")
	in.Command.Flags().StringVarP(&in.Classifier, "classifier", "c", "", "The classifier of the artifact to query")
	in.Command.Flags().StringVarP(&in.Repository, "repository", "r", "https://repo1.maven.org/maven2", "The repository to query")

	return in
}

func (i *InCmd) Run(cmd *cobra.Command, args []string) {
	//fmt.Printf("querying %s:%s:%s from %s\n", i.GroupId, i.ArtifactId, i.Type, i.Repository)

	version, uri, err := download.Download(i.GroupId, i.ArtifactId, ".", i.Repository, "download.jar", "jar", "user", "password")
	if err != nil {
		panic(err)
	}

	out := Out{
		Version: version,
		Uri:     uri,
	}
	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func init() {
	rootCmd.AddCommand(NewInCmd().Command)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
