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
	in := CheckCmd{}
	in.Command = &cobra.Command{
		Use:   "check",
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

func (i *CheckCmd) Run(cmd *cobra.Command, args []string) {
	var jsonIn In

	err := json.NewDecoder(os.Stdin).Decode(&jsonIn)
	if err != nil {
		log.Fatal(err)
	}

	version, _, err := download.Download(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, ".", jsonIn.Source.Repository, "download.jar", jsonIn.Source.Type, "", "")
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

func init() {
	rootCmd.AddCommand(NewCheckCmd().Command)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
