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
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: in.Run,
	}

	return in
}

func (i *InCmd) Run(cmd *cobra.Command, args []string) {
	var jsonIn In

	err := json.NewDecoder(os.Stdin).Decode(&jsonIn)
	if err != nil {
		log.Fatal(err)
	}

	version, uri, err := download.Download(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, ".", jsonIn.Source.Repository, "download.jar", jsonIn.Source.Type, "", "")
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
