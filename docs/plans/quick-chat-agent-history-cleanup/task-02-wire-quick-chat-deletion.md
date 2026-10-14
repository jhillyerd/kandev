---
id: "02-wire-quick-chat-deletion"
title: "Wire cleanup into Quick Chat deletion and expiry"
status: done
wave: 2
depends_on:
  - "01-agent-history-removers"
plan: "plan.md"
requirements:
  - REQ-TASKS-QUICK-CHAT-AGENT-HISTORY-001
  - REQ-TASKS-QUICK-CHAT-AGENT-HISTORY-002
acceptance_criteria:
  - AC-TASKS-QUICK-CHAT-AGENT-HISTORY-001.1
  - AC-TASKS-QUICK-CHAT-AGENT-HISTORY-001.2
  - AC-TASKS-QUICK-CHAT-AGENT-HISTORY-001.4
  - AC-TASKS-QUICK-CHAT-AGENT-HISTORY-001.5
  - AC-TASKS-QUICK-CHAT-AGENT-HISTORY-002.1
  - AC-TASKS-QUICK-CHAT-AGENT-HISTORY-002.4
system_design:
  - ../../specs/tasks/system-design/quick-chat-agent-history-cleanup.md
---

# Task 02: Wire cleanup into Quick Chat deletion and expiry

## Summary

In `performTaskCleanup`, capture each eligible Quick Chat session's history identity before its
executor row is deleted, then call `agenthistory.Remove` after `cleanupQuickChatDirs`. Failures
are logged and never change the deletion result.

## In scope

- Eligibility check (Quick Chat task, not managed conversation, not preserved runtime,
  host-local executor, supported agent, safe token).
- Identity capture: agent name from the session's agent profile snapshot, resume token and
  worktree path from `executors_running`, scratch path from `quickChatDir`, config-dir env from
  the resolved agent profile environment. Confirm the launched cwd for multi-repository chats.
- Injectable remover and home dir on `Service` for tests (follow `SetQuickChatDir`).

## Out of scope

- Archive (does not delete Quick Chats' history).
- Workspace deletion beyond what already routes through `performTaskCleanup`.

## Acceptance

- Service test: create Quick Chat task + session + executor row with a resume token and a fake
  Claude transcript; `DeleteTask` removes only that transcript; a normal task's transcript in the
  same fake root remains after deleting that normal task.
- Service test: `runQuickChatExpiration` with an expired Quick Chat removes its transcript.
- Service test: deleting a Quick Chat with no transcript, no resume token, a non-local executor,
  or a remover error succeeds and leaves the existing deletion outcome unchanged.

## Verification

From `apps/backend`:

```bash
go test -trimpath ./internal/task/service/ -run 'QuickChat' -count=1
go test -trimpath ./internal/task/service/ -count=1
```

From the repository root:

```bash
bash scripts/lint-go-changed $(git diff --name-only main -- '*.go')
```

## Files likely touched

- `apps/backend/internal/task/service/service_tasks.go`
- `apps/backend/internal/task/service/quick_chat_agent_history.go` (new)
- `apps/backend/internal/task/service/quick_chat_agent_history_test.go` (new)
- `apps/backend/internal/backendapp/main.go` (only if the home dir or env needs wiring)

## Dependencies

WO-01.

## Risks

`performTaskCleanup` is shared by every task type; the eligibility check must run first so
normal tasks never reach the remover.

## Inputs

- [Requirements](../../specs/tasks/requirements/quick-chat-agent-history-cleanup.md)
- [System design](../../specs/tasks/system-design/quick-chat-agent-history-cleanup.md)
- `apps/backend/AGENTS.md`

## Results

- Identities are captured in `deleteTaskWithManagedClaim` (non-managed deletes only) and carried
  through `taskEnvironmentCleanup.agentHistory` and the cleanup job snapshot `agent_history`.
- Host-local check uses `executors_running.runtime == standalone`; config dirs come from the
  backend process environment.
- `go test -trimpath ./internal/task/service/... -count=1`: pass (new service tests fail when the
  removal call is disabled). `go vet`, `gofmt -l`: clean. `golangci-lint` v2.9.0
  `--new-from-rev` on both packages: 0 issues. Run in a `golang:1.26` container; Go is not
  installed on the host.
