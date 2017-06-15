# gv

## Summary

Small helper to parse and display build information and glide data in cli and api applications.  This package works with the following github packages: `Masterminds/glide`, `mitchellh/cli` and `emicklei/go-restful`

## Usage

To enable version information to be processed in your application you will need to set and pass a few bits of information in during your build process.  The standard `build.sh` contains the following bits:

```bash
## package declarations
BUILD_NAME="myapp"
BUILD_VERSION="1.0"
BUILD_DATE="$(date -u +%Y%m%d@%H%M%S%z)"
BUILD_BRANCH="$(git branch --no-color|awk '/^\*/ {print $2}')"
BUILD_COMMIT="$(git rev-parse --verify HEAD)"
BUILD_NUMBER="${BAMBOO_BUILD:-00}" ## exported by build system

## encode import data to avoid problems with special characters
GLIDE_DATA=$(base64 glide.yaml | sed ':a;N;$!ba;s/\n//g')
GLIDE_LOCK=$(base64 glide.lock | sed ':a;N;$!ba;s/\n//g')

## build it
go build -o "${BUILD_NAME}" -ldflags "\
-X main.buildName=${BUILD_NAME} \
-X main.buildVersion=${BUILD_VERSION} \
-X main.buildDate=${BUILD_DATE} \
-X main.buildBranch=${BUILD_BRANCH} \
-X main.buildCommit=${BUILD_COMMIT} \
-X main.buildNumber=${BUILD_NUMBER} \
-X main.glideData=${GLIDE_DATA} \
-X main.glideLock=${GLIDE_LOCK}"
```

You then need the following bits in your applications main package:

```go
var (
	// injected build data
	buildName    string
	buildVersion string
	buildDate    string
	buildBranch  string
	buildCommit  string
	buildNumber  string
	glideData    string
	glideLock    string
	// build info struct
	buildInfo *gv.BuildInfo
)

// populate build info
buildInfo = &gv.BuildInfo{
	Name:          buildName,
	Version:       buildVersion,
	Date:          buildDate,
	Branch:        buildBranch,
	Commit:        buildCommit,
	Build:         buildNumber,
	GlideData:     glideData,
	GlideLockData: glideLock,
}
```

To create a `version` command with the `cli` package:

```go
import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/myENA/gv"
	gvCmd "github.com/myENA/gv/cmd"
)


// commands represents all available commands
var commands map[string]cli.CommandFactory

// init command factory
func init() {
	var ui = &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr} // basic ui

	commands = map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return gvCmd.New(buildInfo, ui)
		},
	}
}
```

When you're all done you can see something similar to this in the terminal:

```
$ ./testapp version -d
==>	testapp v1.0
Build:	00
Branch:	master
Commit:	2766668bab9cf87a8fb6331e3dbbf17fdc1f3b78
Date:	20170615@194043+0000

Imports 11 Packages
[001] 4239b770 github.com/armon/go-radix
[002] 4aabc248 github.com/bgentry/speakeasy
[003] f4360770 github.com/hashicorp/consul
[004] 3573b8b5 github.com/hashicorp/go-cleanhttp
[005] 84607742 github.com/Masterminds/glide
[006] 3084677c github.com/Masterminds/vcs
[007] fc9e8d8e github.com/mattn/go-isatty
[008] b481eac7 github.com/mitchellh/cli
[009] 61c576ce github.com/myENA/gv
[010] fb4cac33 golang.org/x/sys
[011] cd8b52f8 gopkg.in/yaml.v2
```

Weee!
