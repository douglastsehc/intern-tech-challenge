package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

type Versions []*semver.Version

func (s Versions) Len() int {
	return len(s)
}

func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Versions) Less(i, j int) bool {
	return s[j].LessThan(*s[i])
}

// Sort sorts the given slice of Version
func Sort(versions []*semver.Version) {
	sort.Sort(Versions(versions))
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version

	//semver.Sort(releases)
	sort.Sort(Versions(releases))
	//fmt.Printf("%s \n", releases)
	for i := range releases {
		var release = releases[i]
		// make sure the version is later than minVersion and without any prerelease
		if release.PreRelease != "" || release.Compare(*minVersion) < 0 {
			continue
		}

		//fmt.Printf("Every version Version: %s   Major: %d    Minor: %d Patch: %d \n\n", release, release.Major, release.Minor, release.Patch)
		if len(versionSlice) == 0 || versionSlice[len(versionSlice)-1].Major != release.Major || versionSlice[len(versionSlice)-1].Minor != release.Minor {
			versionSlice = append(versionSlice, release)
		}
	}
	//fmt.Printf("%s\n", versionSlice)
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, "kubernetes", "kubernetes", opt)
	if err != nil {
		//panic(err) // is this really a good way?
		fmt.Printf("\n%s  \n\n\n", err.Error())
	}
	minVersion := semver.New("1.8.0")
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, minVersion)
	fmt.Printf("latest versions of kubernetes/kubernetes: %s", versionSlice)
}
