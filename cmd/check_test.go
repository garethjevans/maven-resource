package cmd_test

import (
	"testing"

	"github.com/garethjevans/maven-resource/cmd"
	"github.com/garethjevans/maven-resource/download/downloadfakes"
	"github.com/stretchr/testify/assert"
)

func TestCheckCmd_Run_InitialVersion(t *testing.T) {
	fake := &downloadfakes.FakeDownloader{}
	check := cmd.CheckCmd{
		Downloader: fake,
	}

	fake.GetVersionsReturns([]string{"3.0.0", "3.1.0", "3.2.0", "3.3.0"}, nil)

	v, err := check.RunWithInput(cmd.In{
		Source: cmd.Source{
			Repository: "https://repo.spring.io/release",
			ArtifactId: "tomcat-lifecycle-support",
			GroupId:    "org.cloudfoundry",
		},
	})

	t.Logf("versions = %s", v)
	assert.NoError(t, err)
	assert.Equal(t, len(v), 4)
	assert.Equal(t, fake.GetVersionsCallCount(), 1)
}

func TestCheckCmd_Run_AfterVersion(t *testing.T) {
	fake := &downloadfakes.FakeDownloader{}
	check := cmd.CheckCmd{
		Downloader: fake,
	}

	fake.GetVersionsReturns([]string{"3.0.0", "3.1.0", "3.2.0", "3.3.0"}, nil)

	v, err := check.RunWithInput(cmd.In{
		Source: cmd.Source{
			Repository: "https://repo.spring.io/release",
			ArtifactId: "tomcat-lifecycle-support",
			GroupId:    "org.cloudfoundry"},
		Version: cmd.Version{Ref: "3.1.0"},
	})

	t.Logf("versions = %s", v)
	assert.NoError(t, err)
	assert.Equal(t, len(v), 3)
	assert.Equal(t, fake.GetVersionsCallCount(), 1)
}

func TestCheckCmd_Run_AfterVersion_Release(t *testing.T) {
	fake := &downloadfakes.FakeDownloader{}
	check := cmd.CheckCmd{
		Downloader: fake,
	}

	fake.GetVersionsReturns([]string{"3.0.0.RELEASE", "3.1.0.RELEASE", "3.2.0.RELEASE", "3.3.0.RELEASE"}, nil)

	v, err := check.RunWithInput(cmd.In{
		Source: cmd.Source{
			Repository: "https://repo.spring.io/release",
			ArtifactId: "tomcat-lifecycle-support",
			GroupId:    "org.cloudfoundry"},
		Version: cmd.Version{Ref: "3.1.0.RELEASE"},
	})

	t.Logf("versions = %s", v)
	assert.NoError(t, err)
	assert.Equal(t, len(v), 3)
	assert.Equal(t, fake.GetVersionsCallCount(), 1)
	assert.Equal(t, v[0].Ref, "3.1.0.RELEASE")
}
