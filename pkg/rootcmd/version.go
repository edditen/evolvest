package rootcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	buildGitTag    string
	buildGitCommit string
	buildGitBranch string
	buildTime      string
	buildDiff      string
	goVersion      string
	programName    string
)

func init() {
	// get binary location
	file, _ := exec.LookPath(os.Args[0])
	_, programName = filepath.Split(file)
}

// formatVersion format version information to string.
func formatVersion() string {
	stringBuilder := strings.Builder{}
	f := func(s string, args ...interface{}) {
		stringBuilder.WriteString(fmt.Sprintf(s, args...))
	}
	f("\n--------------\n")
	f("[%s] tag: %s\n\n", programName, buildGitTag)
	f("Git commit: %s\n", buildGitCommit)
	f("Git branch: %s\n", buildGitBranch)
	f("Build time: %s\n", buildTime)
	f("Go version: %s\n", goVersion)

	if buildDiff == "" {
		f(`
WARN: No version found, properly because the codebase wasn't under version control.
      Or the binary wasn't build with correct FLAGS.
`)
	} else if buildDiff != "0" {
		f("\nWARN: %s changes not commit\n", buildDiff)
	}
	f("--------------\n")

	return stringBuilder.String()
}
