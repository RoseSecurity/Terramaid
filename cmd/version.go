package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var Version = "1.6.2"

type Release struct {
	TagName string `json:"tag_name"`
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the CLI version",
	Long:    `This command prints the CLI version`,
	Example: "terramaid version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("terramaid: " + Version)
		latestReleaseTag, err := latestRelease()
		if err == nil && latestReleaseTag != "" {
			latestRelease := strings.TrimPrefix(latestReleaseTag, "v")
			currentRelease := strings.TrimPrefix(Version, "v")
			if latestRelease != currentRelease {
				updateTerramaid(latestRelease)
			}
		}
	},
}

func latestRelease() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/RoseSecurity/terramaid/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch version: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to parse version: %w", err)
	}

	return release.TagName, nil
}

func updateTerramaid(latestVersion string) {
	c1 := color.New(color.FgCyan)

	c1.Println(fmt.Sprintf("\nYour version of Terramaid is out of date. The latest version is %s\n\n", latestVersion))
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
