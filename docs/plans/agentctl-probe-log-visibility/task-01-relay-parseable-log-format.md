---
id: "01-relay-parseable-log-format"
title: "Pin a relay-parseable agentctl log format"
status: done
wave: 1
depends_on: []
plan: "plan.md"
requirements:
  - REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002
acceptance_criteria:
  - AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002.1
system_design:
  - ../../specs/agents/system-design/probe-failure-diagnostics.md
---

# Task 01: Pin a relay-parseable agentctl log format

## Summary

The launcher owns the agentctl child's log format because the child's stdout
feeds the launcher's relay. Pin `AGENTCTL_LOG_FORMAT` to the console format the
relay parses and drop any inherited value, so an operator's ambient environment
cannot break lifecycle record ingestion.

## In scope

- Add `AGENTCTL_LOG_FORMAT=console` to `Launcher.childEnvOverrides` in
  `internal/agent/runtime/agentctl/launcher/launcher.go`.
- Regression test proving the override pins the format and replaces an
  inherited value.

## Out of scope

- The relay's record parsing (landed upstream in #4223 and #4263).
- Any probe-side behavior (Task 02).

## Applicable criteria

- AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002.1
- System design: [Probe Failure Diagnostics](../../specs/agents/system-design/probe-failure-diagnostics.md)

## Acceptance conditions

1. `TestChildEnvOverridesPinRelayParseableLogFormat` fails before the pin and
   passes after.
2. An inherited `AGENTCTL_LOG_FORMAT` value cannot survive into the child
   environment.

## Verification commands

```bash
cd apps/backend
go test ./internal/agent/runtime/agentctl/launcher/ -run TestChildEnvOverridesPinRelayParseableLogFormat -count=1
```

## Likely files

- `apps/backend/internal/agent/runtime/agentctl/launcher/launcher.go`
- `apps/backend/internal/agent/runtime/agentctl/launcher/launcher_test.go`

## Results

Implemented in commit `33980ff74`; the regression test passes and the
launcher package tests are green.
