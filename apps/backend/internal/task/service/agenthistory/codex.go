package agenthistory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// codexMetaLineLimit bounds how much of a rollout's first line is read.
const codexMetaLineLimit = 1 << 20

type codexSessionMeta struct {
	Type    string `json:"type"`
	Payload struct {
		ID  string `json:"id"`
		Cwd string `json:"cwd"`
	} `json:"payload"`
}

// removeCodex removes <root>/sessions/YYYY/MM/DD/rollout-*-<token>.jsonl files
// whose session_meta records the same thread ID and working directory. Codex's
// shared thread indexes are left to Codex, which skips threads without a rollout.
func removeCodex(ctx context.Context, id Identity) error {
	root := configRoot(id, "CODEX_HOME", ".codex")
	if root == "" {
		return nil
	}
	sessions := filepath.Join(root, "sessions")
	if ok, err := realDir(sessions); err != nil || !ok {
		return err
	}
	matches, err := filepath.Glob(filepath.Join(sessions, "*", "*", "*", "rollout-*-"+id.ResumeToken+".jsonl"))
	if err != nil {
		return err
	}
	workDir := filepath.Clean(id.WorkDir)
	var errs []error
	for _, path := range matches {
		if err := context.Cause(ctx); err != nil {
			return errors.Join(append(errs, err)...)
		}
		owned, err := codexRolloutBelongsTo(path, id.ResumeToken, workDir)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if owned {
			errs = append(errs, removeRegularFile(path))
		}
	}
	return errors.Join(errs...)
}

func codexRolloutBelongsTo(path, thread, workDir string) (bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil || !info.Mode().IsRegular() {
		return false, err
	}
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), codexMetaLineLimit)
	if !scanner.Scan() {
		return false, scanner.Err()
	}
	var meta codexSessionMeta
	if json.Unmarshal(scanner.Bytes(), &meta) != nil {
		return false, nil
	}
	return meta.Type == "session_meta" &&
		meta.Payload.ID == thread &&
		meta.Payload.Cwd != "" &&
		filepath.Clean(meta.Payload.Cwd) == workDir, nil
}
