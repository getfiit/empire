package api_test

import (
	"testing"

	"github.com/getfiit/empire/pkg/heroku"
)

func TestReleaseList(t *testing.T) {
	c, s := NewTestClient(t)
	defer s.Close()

	mustDeploy(t, c, DefaultImage)

	releases := mustReleaseList(t, c, "acme-inc")

	if len(releases) != 1 {
		t.Fatal("Expected a release")
	}

	if got, want := releases[0].Version, 1; got != want {
		t.Fatalf("Version => %v; want %v", got, want)
	}
}

func TestReleaseInfo(t *testing.T) {
	c, s := NewTestClient(t)
	defer s.Close()

	mustDeploy(t, c, DefaultImage)

	release := mustReleaseInfo(t, c, "acme-inc", "1")

	if got, want := release.Version, 1; got != want {
		t.Fatalf("Version => %v; want %v", got, want)
	}
}

func TestReleaseRollback(t *testing.T) {
	c, s := NewTestClient(t)
	defer s.Close()

	// Deploy twice
	mustDeploy(t, c, DefaultImage)
	mustDeploy(t, c, DefaultImage)

	// Rollback to the first deploy.
	mustReleaseRollback(t, c, "acme-inc", "1")
}

func mustReleaseList(t testing.TB, c *heroku.Client, appName string) []heroku.Release {
	releases, err := c.ReleaseList(appName, nil)
	if err != nil {
		t.Fatal(err)
	}

	return releases
}

func mustReleaseInfo(t testing.TB, c *heroku.Client, appName string, version string) *heroku.Release {
	release, err := c.ReleaseInfo(appName, version)
	if err != nil {
		t.Fatal(err)
	}

	return release
}

func mustReleaseRollback(t testing.TB, c *heroku.Client, appName string, version string) *heroku.Release {
	release, err := c.ReleaseRollback(appName, version, "")
	if err != nil {
		t.Fatal(err)
	}

	return release
}
