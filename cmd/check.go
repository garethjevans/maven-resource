package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/garethjevans/maven-resource/download"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var semverRE = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

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

	versionToCheck, err := semver.NewVersion(jsonIn.Version.Ref)
	if err != nil {
		log.Fatalf("Error parsing version: %s: %s", jsonIn.Version.Ref, err)
	}

	versions, err := download.GetVersions(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, jsonIn.Source.Repository, jsonIn.Source.Username, jsonIn.Source.Password)
	if err != nil {
		panic(err)
	}

	var refs []Version
	for _, versionString := range versions {
		if semverRE.MatchString(versionString) {
			v, err := semver.NewVersion(versionString)
			if err != nil {
				log.Fatalf("Error parsing version: %s: %s", versionString, err)
			}

			if versionToCheck.LessThan(v) {
				refs = append(refs, Version{Ref: versionString})
			}
		}
	}
	b, err := json.Marshal(refs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
