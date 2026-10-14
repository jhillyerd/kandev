---
id: "02-probe-failure-log-diagnostics"
title: "Log sanitized probe failure diagnostics"
status: done
wave: 1
depends_on: []
plan: "plan.md"
requirements:
  - REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001
acceptance_criteria:
  - AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.1
  - AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.2
  - AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.3
  - AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.4
system_design:
  - ../../specs/agents/system-design/probe-failure-diagnostics.md
---

# Task 02: Log sanitized probe failure diagnostics

## Summary

A failed probe leaves child-side failure evidence in the backend log: a
bounded stderr tail, sanitized for the log boundary, plus the failure
classification the client response cannot carry. Client payloads stay
sanitized.

## In scope

- Capture the child's stderr in the trailing-window `stderrBuffer`
  (4 KiB retention) during probes.
- On failure, log the tail at both failure sites (`ACP probe failed`,
  `ACP profile probe failed`) through `probeStderrForLog`
  (`routingerr.SanitizeCredentials` composed with `scrubTmpPaths`), so
  credential-shaped values are masked and tmp paths replaced while diagnostic
  text stays intact.
- In profile context, additionally log the pre-sanitization managed-runtime
  classification as `managed_failure_code` alongside the client-visible code.
- Regression tests for the scrubbed tail, the credential masking (both
  contexts), and the pre-sanitize code.

## Out of scope

- Client-visible probe error copy and provider-gateway status hints.
- Persistent-session stderr streaming to the task UI.
- The launcher log-format pin (Task 01).

## Applicable criteria

- AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.1 through 001.4
- System design: [Probe Failure Diagnostics](../../specs/agents/system-design/probe-failure-diagnostics.md)

## Acceptance conditions

1. A profile probe failure log entry carries the npm/child cause with tmp
   paths scrubbed and the managed failure code, while the client response
   keeps the sanitized message and empty code.
2. Credential-shaped values in child stderr (token prefixes, Bearer
   challenges, credential assignments) are masked in both failure logs, and
   the diagnostic text around them survives.
3. The captured tail never exceeds the trailing window.

## Verification commands

```bash
cd apps/backend
go test ./internal/agentctl/server/utility/ -run 'ProbeFailureLog|CarriesScrubbedStderr|MasksCredentialShapedStderr' -count=1
```

## Likely files

- `apps/backend/internal/agentctl/server/utility/acp_executor.go`
- `apps/backend/internal/agentctl/server/utility/profile_probe_context_test.go`
- `apps/backend/internal/agentctl/server/utility/acp_executor_test.go`

## Results

Implemented across commits `7f16a1cd8` (scrubbed tail and managed failure
code) and `09a0eca40` (credential masking at both log sites, from PR review).
The bounded stderr capture and its tests already existed in the probe
executor. All listed tests pass; the utility package is green apart from two
pre-existing NixOS-environment failures confirmed identical at baseline.
