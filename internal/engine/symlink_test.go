package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSymlink(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source")
	dst := filepath.Join(dir, "link")
	require.NoError(t, os.WriteFile(src, []byte("hello"), 0o644))

	err := CreateSymlink(src, dst)
	require.NoError(t, err)

	isLink, err := IsSymlink(dst)
	require.NoError(t, err)
	assert.True(t, isLink)

	target, _ := os.Readlink(dst)
	assert.Equal(t, src, target)
}

func TestCreateSymlinkOverwrite(t *testing.T) {
	dir := t.TempDir()
	src1 := filepath.Join(dir, "source1")
	src2 := filepath.Join(dir, "source2")
	dst := filepath.Join(dir, "link")
	os.WriteFile(src1, []byte("one"), 0o644)
	os.WriteFile(src2, []byte("two"), 0o644)

	CreateSymlink(src1, dst)
	CreateSymlink(src2, dst)

	target, _ := os.Readlink(dst)
	assert.Equal(t, src2, target)
}

func TestRemoveSymlink(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source")
	dst := filepath.Join(dir, "link")
	os.WriteFile(src, []byte("hello"), 0o644)
	CreateSymlink(src, dst)

	err := RemoveSymlink(dst)
	require.NoError(t, err)

	_, err = os.Lstat(dst)
	assert.True(t, os.IsNotExist(err))
}

func TestRemoveGeneratedContents(t *testing.T) {
	dir := t.TempDir()
	genDir := filepath.Join(dir, ".generated")
	os.MkdirAll(filepath.Join(genDir, "claude-code"), 0o755)
	os.WriteFile(filepath.Join(genDir, "claude-code", "test.md"), []byte("test"), 0o644)

	err := RemoveGeneratedContents(dir)
	require.NoError(t, err)

	entries, _ := os.ReadDir(genDir)
	assert.Equal(t, 0, len(entries))
}

func TestRemoveStaleSymlinks(t *testing.T) {
	dir := t.TempDir()
	dotgenDir := filepath.Join(dir, ".dotagen")
	genDir := filepath.Join(dotgenDir, ".generated", "claude-code")
	os.MkdirAll(genDir, 0o755)

	activeFile := filepath.Join(genDir, "da-active.md")
	staleFile := filepath.Join(genDir, "da-stale.md")
	os.WriteFile(activeFile, []byte("active"), 0o644)
	os.WriteFile(staleFile, []byte("stale"), 0o644)

	claudeDir := filepath.Join(dir, ".claude", "agents")
	os.MkdirAll(claudeDir, 0o755)
	CreateSymlink(activeFile, filepath.Join(claudeDir, "da-active.md"))
	CreateSymlink(staleFile, filepath.Join(claudeDir, "da-stale.md"))

	removed, err := RemoveStaleSymlinks(dir, dotgenDir, []string{"da-active"}, []string{"claude-code"})
	require.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(".claude", "agents", "da-stale.md")}, removed)

	_, err = os.Lstat(filepath.Join(claudeDir, "da-active.md"))
	assert.NoError(t, err)
	_, err = os.Lstat(filepath.Join(claudeDir, "da-stale.md"))
	assert.True(t, os.IsNotExist(err))
}

func TestFindDotagenSymlinksSkipsNonDaAgents(t *testing.T) {
	dir := t.TempDir()
	dotgenDir := filepath.Join(dir, ".dotagen")
	genDir := filepath.Join(dotgenDir, ".generated", "claude-code")
	os.MkdirAll(genDir, 0o755)

	daFile := filepath.Join(genDir, "da-myagent.md")
	userFile := filepath.Join(genDir, "my-custom-agent.md")
	os.WriteFile(daFile, []byte("da"), 0o644)
	os.WriteFile(userFile, []byte("user"), 0o644)

	claudeDir := filepath.Join(dir, ".claude", "agents")
	os.MkdirAll(claudeDir, 0o755)
	CreateSymlink(daFile, filepath.Join(claudeDir, "da-myagent.md"))
	CreateSymlink(userFile, filepath.Join(claudeDir, "my-custom-agent.md"))

	links, err := FindDotagenSymlinks(dir, dotgenDir)
	require.NoError(t, err)

	assert.Len(t, links, 1)
	assert.Equal(t, "da-myagent", links[0].Agent)
}
