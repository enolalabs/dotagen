package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/enolalabs/dotagen/v2/skillsrc"
)

func TestRunInitCopiesHallmarkAndLeavesItDisabled(t *testing.T) {
	home := t.TempDir()
	if _, _, _, err := runInit(home); err != nil {
		t.Fatal(err)
	}

	const name = "dotagent-nutlope-hallmark"
	embeddedFiles := skillsrc.ListSkillFiles(name)
	if len(embeddedFiles) == 0 {
		t.Fatalf("no embedded files for %s — skill not bundled", name)
	}
	sort.Strings(embeddedFiles)

	for _, rel := range embeddedFiles {
		want, err := skillsrc.ReadSkillFile(name + "/" + rel)
		if err != nil {
			t.Fatalf("embedded %s: %v", rel, err)
		}
		got, err := os.ReadFile(filepath.Join(home, ".dotagen", "skills", name, rel))
		if err != nil {
			t.Fatalf("initialized %s: %v", rel, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("%s differs from embedded content", rel)
		}
	}

	skillDir := filepath.Join(home, ".dotagen", "skills", name)
	var destFiles []string
	err := filepath.Walk(skillDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(skillDir, path)
		if err != nil {
			return err
		}
		destFiles = append(destFiles, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", skillDir, err)
	}
	sort.Strings(destFiles)

	if !reflect.DeepEqual(destFiles, embeddedFiles) {
		t.Fatalf("destination files != embedded files\ndest:     %v\nembedded: %v", destFiles, embeddedFiles)
	}

	config, err := os.ReadFile(filepath.Join(home, ".dotagen", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(config), "  "+name+":\n    targets: []\n") {
		t.Fatalf("Hallmark is not disabled: %s", config)
	}
}
