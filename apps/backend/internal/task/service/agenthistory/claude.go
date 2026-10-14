package agenthistory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// claudeMaxProjectNameLength is the longest encoded project name Claude Code
// stores verbatim; longer names are truncated and hashed by Claude.
const claudeMaxProjectNameLength = 200

const claudeMemoryDir = "memory"

// encodeClaudeProjectName mirrors Claude Code's project directory naming:
// every character outside [A-Za-z0-9] becomes '-'.
func encodeClaudeProjectName(workDir string) string {
	var b strings.Builder
	b.Grow(len(workDir))
	for _, r := range workDir {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		b.WriteByte('-')
	}
	return b.String()
}

// removeClaude removes <root>/projects/<encoded cwd>/<token>.jsonl, its
// <token>/ sidecar directory, and the project directory when nothing else
// remains in it.
func removeClaude(ctx context.Context, id Identity) error {
	root := configRoot(id, "CLAUDE_CONFIG_DIR", ".claude")
	if root == "" {
		return nil
	}
	name := encodeClaudeProjectName(filepath.Clean(id.WorkDir))
	if len(name) > claudeMaxProjectNameLength {
		return nil
	}
	projects := filepath.Join(root, "projects")
	if ok, err := realDir(projects); err != nil || !ok {
		return err
	}
	project := filepath.Join(projects, name)
	if ok, err := realDir(project); err != nil || !ok {
		return err
	}
	if err := context.Cause(ctx); err != nil {
		return err
	}
	var errs []error
	if err := removeRegularFile(filepath.Join(project, id.ResumeToken+".jsonl")); err != nil {
		errs = append(errs, err)
	}
	sidecar := filepath.Join(project, id.ResumeToken)
	if ok, err := realDir(sidecar); err != nil {
		errs = append(errs, err)
	} else if ok {
		if err := os.RemoveAll(sidecar); err != nil {
			errs = append(errs, fmt.Errorf("remove %s: %w", sidecar, err))
		}
	}
	if len(errs) == 0 {
		errs = append(errs, removeEmptyClaudeProject(project))
	}
	return errors.Join(errs...)
}

// removeEmptyClaudeProject removes a project directory that holds nothing, or
// only an empty memory directory.
func removeEmptyClaudeProject(project string) error {
	entries, err := os.ReadDir(project)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() != claudeMemoryDir || !entry.IsDir() {
			return nil
		}
		memory, err := os.ReadDir(filepath.Join(project, claudeMemoryDir))
		if err != nil || len(memory) > 0 {
			return err
		}
	}
	if len(entries) > 0 {
		if err := os.Remove(filepath.Join(project, claudeMemoryDir)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err := os.Remove(project); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", project, err)
	}
	return nil
}
