package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/kandev/kandev/internal/agentruntime"
	"github.com/kandev/kandev/internal/task/models"
	"github.com/kandev/kandev/internal/task/service/agenthistory"
)

// agentHistoryEnvKeys are the agent config-dir overrides read from the backend
// environment, which host-local agent processes inherit.
var agentHistoryEnvKeys = []string{"CLAUDE_CONFIG_DIR", "CODEX_HOME"}

// agentHistoryTarget is the durable identity of one Quick Chat session's agent
// history, captured before the executor row that records it is deleted.
type agentHistoryTarget struct {
	SessionID   string            `json:"session_id"`
	AgentName   string            `json:"agent_name"`
	ResumeToken string            `json:"resume_token"`
	WorkDir     string            `json:"work_dir"`
	HomeDir     string            `json:"home_dir,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
}

// captureQuickChatAgentHistory returns the agent history targets of a Quick
// Chat's host-local sessions. Any other task yields none.
func (s *Service) captureQuickChatAgentHistory(
	ctx context.Context, task *models.Task, sessions []*models.TaskSession,
) []agentHistoryTarget {
	if !IsRestorableQuickChatTask(task) || s.executors == nil {
		return nil
	}
	home, _ := os.UserHomeDir()
	env := make(map[string]string, len(agentHistoryEnvKeys))
	for _, key := range agentHistoryEnvKeys {
		if value := os.Getenv(key); value != "" {
			env[key] = value
		}
	}
	var targets []agentHistoryTarget
	for _, session := range sessions {
		if session == nil || session.ID == "" {
			continue
		}
		running, err := s.executors.GetExecutorRunningBySessionID(ctx, session.ID)
		if err != nil {
			if !errors.Is(err, models.ErrExecutorRunningNotFound) {
				s.logger.Warn("quick-chat agent history: failed to read executor runtime",
					zap.String("task_id", task.ID), zap.String("session_id", session.ID), zap.Error(err))
			}
			continue
		}
		if running == nil || running.Runtime != agentruntime.RuntimeStandalone || running.ResumeToken == "" {
			continue
		}
		workDir := running.WorktreePath
		if workDir == "" && s.quickChatDir != "" {
			workDir = filepath.Join(s.quickChatDir, session.ID)
		}
		targets = append(targets, agentHistoryTarget{
			SessionID:   session.ID,
			AgentName:   models.StringFromAny(session.AgentProfileSnapshot["agent_name"]),
			ResumeToken: running.ResumeToken,
			WorkDir:     workDir,
			HomeDir:     home,
			Env:         env,
		})
	}
	return targets
}

// removeQuickChatAgentHistory removes captured agent history. Failures are
// logged and never change the deletion outcome.
func (s *Service) removeQuickChatAgentHistory(
	ctx context.Context, taskID string, targets []agentHistoryTarget, preserveExecutorRows map[string]struct{},
) {
	for _, target := range targets {
		if _, preserve := preserveExecutorRows[target.SessionID]; preserve {
			continue
		}
		err := agenthistory.Remove(ctx, agenthistory.Identity{
			AgentName:   target.AgentName,
			ResumeToken: target.ResumeToken,
			WorkDir:     target.WorkDir,
			Env:         target.Env,
			HomeDir:     target.HomeDir,
		})
		if err != nil {
			s.logger.Warn("failed to remove quick-chat agent history",
				zap.String("task_id", taskID),
				zap.String("session_id", target.SessionID),
				zap.String("agent", target.AgentName),
				zap.Error(err))
		}
	}
}
