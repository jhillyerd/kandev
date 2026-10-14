package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestDeleteQuickChatRemovesHistoryBeforeResumeTokenRotation(t *testing.T) {
	for _, nextToken := range []string{"", "abababab-0000-4000-8000-000000000006"} {
		t.Run("next="+nextToken, func(t *testing.T) {
			svc, _, repo := createTestService(t)
			f := newAgentHistoryFixture(t, svc)
			token := "abababab-0000-4000-8000-000000000005"
			sessionID := seedAgentHistoryTask(t, repo, agentHistoryTaskSpec{
				taskID: "quick-rotated-history", ephemeral: true, agentName: "claude-acp", resumeToken: token,
			})
			workDir := filepath.Join(f.quickChatDir, sessionID)
			before := f.writeClaudeTranscript(t, workDir, token)
			if err := repo.UpdateResumeToken(context.Background(), sessionID, "", nextToken, ""); err != nil {
				t.Fatalf("UpdateResumeToken: %v", err)
			}
			var after string
			if nextToken != "" {
				after = f.writeClaudeTranscript(t, workDir, nextToken)
			}
			if err := svc.DeleteTask(context.Background(), "quick-rotated-history"); err != nil {
				t.Fatalf("DeleteTask: %v", err)
			}
			waitForCleanupDone(t, svc)
			assertPathExists(t, before, false)
			if after != "" {
				assertPathExists(t, after, false)
			}
		})
	}
}

func TestDeleteQuickChatFailedStopPreservesLiveHistory(t *testing.T) {
	svc, _, repo := createTestService(t)
	f := newAgentHistoryFixture(t, svc)
	token := "abababab-0000-4000-8000-000000000007"
	sessionID := seedAgentHistoryTask(t, repo, agentHistoryTaskSpec{
		taskID: "quick-live-history", ephemeral: true, agentName: "claude-acp", resumeToken: token,
	})
	transcript := f.writeClaudeTranscript(t, filepath.Join(f.quickChatDir, sessionID), token)
	stopped := make(chan struct{}, 1)
	svc.executionStopper = &stubStopper{stopSessionErr: errors.New("stop blocked"), onStopSession: func() {
		select {
		case stopped <- struct{}{}:
		default:
		}
	}}
	if err := svc.DeleteTask(context.Background(), "quick-live-history"); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
	waitForCleanupDone(t, svc)
	select {
	case <-stopped:
	default:
		t.Fatal("the failing stop boundary was not exercised")
	}
	assertPathExists(t, transcript, true)
}
