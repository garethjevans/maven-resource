package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/garethjevans/maven-resource/download"

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

	Log("Checking resource for %+v\n", jsonIn)

	var versionToCheck *semver.Version
	if jsonIn.Version.Ref != "" {
		versionToCheck, err = semver.NewVersion(jsonIn.Version.Ref)
		if err != nil {
			Log("Skipping existing version %+s, %s\n", jsonIn.Version.Ref, err)
		}
	}

	versions, err := download.GetVersions(jsonIn.Source.GroupId, jsonIn.Source.ArtifactId, jsonIn.Source.Repository, jsonIn.Source.Username, jsonIn.Source.Password)
	if err != nil {
		panic(err)
	}

	fmt.Println("versions: ", versions)

	var refs []Version
	for _, versionString := range versions {
		compVersionString := versionString
		versionString = strings.Replace(versionString, ".RELEASE", "", -1)

		if semverRE.MatchString(versionString) {
			v, err := semver.NewVersion(versionString)
			if err != nil {
				log.Fatalf("Should never happen! Error parsing version: %s: %s", versionString, err)
			}

			if strings.Contains(compVersionString, ".RELEASE") {
				versionString = versionString + ".RELEASE"
			}

			if versionToCheck == nil || versionToCheck.LessThan(v) {
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
