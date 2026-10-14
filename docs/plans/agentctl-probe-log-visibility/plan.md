---
created: 2026-10-14
status: complete
requirements:
  - REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001
  - REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002
system_design:
  - ../../specs/agents/system-design/probe-failure-diagnostics.md
legacy_specs: []
---

# Implementation Plan: Fix agentctl probe log visibility

## Overview

Make failed capability and profile probes diagnosable from backend logs. Two
work orders restore the operator's window into child-side probe failures and
protect the record ingestion path that depends on the agentctl child's stdout
format.

The agent system owns this repair because it owns the agent capability surface:
the probe executor and the agentctl child launcher are agent-system
components, and the failed outcome is an operator who cannot diagnose a probe
failure. The intended behavior was missing from requirements, so this package
adds the smallest requirement set that defines it
([Probe Failure Diagnostics](../../specs/agents/requirements/probe-failure-diagnostics.md));
no existing requirement or acceptance criterion covers probe failure logging or
the child log format contract.

## Evidence and root cause

A probe whose child dies before answering reported only "peer disconnected
before response". The child's actual reason, such as an `npx` resolution
failure, a missing runtime, or a panic, went to the void: the UI's "check agent
logs" advice pointed at logs that never carried it. Separately, the agentctl
child inherits its environment from the operator's shell, so an ambient
`AGENTCTL_LOG_FORMAT` value could switch the child's stdout to a format the
launcher relay does not parse, silently breaking lifecycle record ingestion.

Root causes:

- The probe executor discarded child stderr instead of capturing a bounded
  tail for the failure log.
- The launcher did not own the child's log format even though the child's
  stdout feeds the launcher's relay.

The relay's parsing of agentctl JSON log records landed separately upstream
(#4223, #4263) and needed no change.

## Work orders

- [Task 01: Pin a relay-parseable agentctl log format](task-01-relay-parseable-log-format.md)
- [Task 02: Log sanitized probe failure diagnostics](task-02-probe-failure-log-diagnostics.md)

Both work orders are independent (different files, different contracts) and
run in wave 1.

## Verification strategy

Each work order owns its regression tests at the owning boundary; there is no
UI change, so no rendered checks apply. Commands:

```bash
cd apps/backend
go test ./internal/agent/runtime/agentctl/launcher/ -run TestChildEnvOverridesPinRelayParseableLogFormat -count=1
go test ./internal/agentctl/server/utility/ -run 'ProbeFailureLog|CarriesScrubbedStderr' -count=1
```

## Risks and exclusions

- Log sanitization must not destroy the diagnostic text the log exists for:
  the credential-only redaction tier is chosen for that reason.
- Client responses stay sanitized; no subprocess detail crosses the client
  boundary.
- Out of scope: the relay's record parsing (upstream), client-visible probe
  error copy, and persistent-session stderr streaming.

## Results

Both work orders are complete and merged in this pull request. The
documentation coverage check is satisfied by this package.
