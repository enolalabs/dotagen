package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const githubRepo = "enolalabs/dotagen"

type ghRelease struct {
	TagName string `json:"tag_name"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dotagen to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVer := normalizeVersion(version)

		fmt.Print("Checking latest version... ")
		latest, err := getLatestVersion()
		if err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("failed to check latest version: %w", err)
		}

		fmt.Printf("%s\n", latest)

		if latest == "v"+currentVer {
			fmt.Printf("Already up to date (v%s).\n", currentVer)
			return nil
		}

		if currentVer != "dev" {
			fmt.Printf("Updating: v%s → %s\n", currentVer, latest)
		} else {
			fmt.Printf("Installing: %s\n", latest)
		}

		plat := runtime.GOOS
		arch := runtime.GOARCH

		assetName := fmt.Sprintf("dotagen_%s_%s_%s", strings.TrimPrefix(latest, "v"), plat, arch)
		ext := ""
		if plat == "windows" {
			assetName += ".exe"
			ext = ".exe"
		}

		url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", githubRepo, latest, assetName)
		fmt.Printf("Downloading %s/%s... ", plat, arch)

		binaryPath, _ := os.Executable()
		if binaryPath == "" {
			binaryPath = "dotagen" + ext
		}

		tmpFile := binaryPath + ".new"
		if err := downloadFile(url, tmpFile); err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("download failed: %w", err)
		}
		os.Chmod(tmpFile, 0o755)
		fmt.Println("OK")

		oldFile := binaryPath + ".old"
		os.Remove(oldFile)
		os.Rename(binaryPath, oldFile)
		if err := os.Rename(tmpFile, binaryPath); err != nil {
			os.Rename(oldFile, binaryPath)
			return fmt.Errorf("failed to replace binary: %w", err)
		}
		os.Remove(oldFile)

		out, _ := exec.Command(binaryPath, "version").Output()
		newVer := strings.TrimSpace(string(out))
		if newVer != "" {
			fmt.Printf("\n%s\n", newVer)
		}
		fmt.Printf("\ndotagen updated successfully!\n")
		return nil
	},
}

func getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	return rel.TagName, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
