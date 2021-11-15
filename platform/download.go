package platform

import (
	"fmt"
	"log"
	"os/user"
	"runtime"

	"github.com/cavaliercoder/grab"
)

type PlatformBundle struct {
	platform string
	version  string
}

func (platformBundle *PlatformBundle) SetPlatform(platform string) {
	platformBundle.platform = platform
}

func (platformBundle *PlatformBundle) SetVersion(version string) {
	platformBundle.version = version
}

func getZipName(platformBundle *PlatformBundle) string {
	return platformBundle.platform + "-" + platformBundle.version + "-" + getOsPlatform() + "-x86_64.zip"
}

func (platformBundle *PlatformBundle) getUrl() string {
	return "https://artifactory.siren.io/artifactory/libs-release-staging-local/solutions/siren/" + platformBundle.platform + "/" + platformBundle.version + "/" + getZipName(platformBundle)
}

func getOsPlatform() string {
	return runtime.GOOS
}

func Download(platform string, version string) (string, string) {

	platform_bundle := PlatformBundle{platform, version}

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	downloadDirectory := user.HomeDir + "/Downloads"

	resp, err := grab.Get(downloadDirectory, platform_bundle.getUrl())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download saved to", resp.Filename)

	return downloadDirectory, resp.Filename
}
