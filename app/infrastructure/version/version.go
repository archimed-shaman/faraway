package version

import "fmt"

var (
	versionExt   string
	buildDateExt string
	gitBranch    string
	gitHash      string
	gitDirty     string
)

func isDirty(c string) string {
	if c == "1" || c == "true" {
		return "dirty"
	}

	return ""
}

func GetProductName() string { return "World of Wisdom" }
func GetVersion() string     { return versionExt }
func GetDate() string        { return buildDateExt }
func GetGit() string         { return fmt.Sprintf("%s %s (%s)", gitBranch, gitHash, isDirty(gitDirty)) }
