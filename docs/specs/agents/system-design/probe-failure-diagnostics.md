---
status: current
system: agents
created: 2026-10-14
updated: 2026-10-14
requirements:
  - REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001
  - REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002
owners:
  - kandev
---
# Probe Failure Diagnostics System Design

## Purpose and boundaries

This design owns the technical contract for
REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001 and
REQ-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002: how a one-shot probe subprocess's
stderr is captured, classified, sanitized, and written to the backend log, and
how the agentctl child's log format is pinned by the launcher that reads its
stdout.

It does not change client-visible probe error copy (owned by
[OpenAI-compatible providers](openai-compatible-providers.md)
for gateway status surfacing), probe command resolution (owned by
[Operator-Registered Agents](custom-acp-agents.md)),
or the launcher relay's record parsing, which is unchanged upstream behavior.

## Requirement mapping

| Acceptance criterion | Owner |
| --- | --- |
| AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.1 | `internal/agentctl/server/utility` probe failure logging |
| AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.2 | `internal/agentctl/server/utility` (`probeStderrForLog`) |
| AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.3 | `internal/agentctl/server/utility` (`sanitizeProfileProbeFailure`) |
| AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-001.4 | `internal/agentctl/server/utility` (`stderrBuffer`) |
| AC-AGENTS-PROBE-FAILURE-DIAGNOSTICS-002.1 | `internal/agent/runtime/agentctl/launcher` (`childEnvOverrides`) |

## Probe stderr capture

The probe executor attaches a trailing-window buffer to the child's stderr.
The buffer keeps only the last `stderrTailLimit` (4 KiB) bytes: `Write` discards
overflow from the front, so a chatty or looping child cannot grow the capture
or the log entry that carries it. An empty tail is itself a result; it means
the child produced no diagnostics, not that diagnostics were lost.

## Classify on raw stderr, sanitize for the log

Failure classification consumes the raw tail, because marker matching
(`managedRuntimeProbeFailureCode`, the authentication markers of
`sanitizeProfileProbeFailure`) and the upstream status extraction of
`withUpstreamHint` are defined over the child's actual output. The upstream
hint only extracts an HTTP-status-shaped match, so no credential can reach a
client message through it.

The log boundary applies the opposite rule. `probeStderrForLog` composes
`routingerr.SanitizeCredentials` with `scrubTmpPaths`:

- The credential-only redaction tier masks credential-shaped values (token
  prefixes such as `sk-`, `ghp_`, `github_pat_`, `kandev_pat_`, `Bearer`
  challenges, `Authorization` headers, `--api-key` values, credential
  assignments, URL userinfo) while leaving error text, paths, and opaque
  identifiers intact, preserving the diagnostic value the log exists for.
- The tmp-path scrub then replaces local temporary path shapes with `<path>`,
  keeping subprocess filesystem details out of backend logs.

Both probe failure sites (`ACP probe failed` and `ACP profile probe failed`)
log the sanitized tail. The profile site additionally logs the client-visible
failure code and the pre-sanitization managed-runtime classification
(`managed_failure_code`), because profile-context responses deliberately
replace both with a generic sanitized message and code.

The client response never carries the tail. Handlers keep subprocess
diagnostics and tmp paths out of client payloads; the log is the operator's
only window into the child-side cause.

## Launcher-owned child log format

The agentctl child's stdout feeds the supervising launcher's relay, never a log
collector. The relay parses a fixed record format, so the log format is the
launcher's contract: `childEnvOverrides` pins `AGENTCTL_LOG_FORMAT` to the
console format the relay parses. Because the pin is an environment override,
any inherited value of that variable from the operator's ambient environment
is dropped rather than merged.

## Verification

Regression tests cover each criterion at the owning boundary:

- `TestChildEnvOverridesPinRelayParseableLogFormat`
  (`internal/agent/runtime/agentctl/launcher`)
- `TestProfileProbeFailureLogCarriesScrubbedStderr` and
  `TestProfileProbeFailureLogMasksCredentialShapedStderr`
  (`internal/agentctl/server/utility`, profile context)
- `TestProbeFailureLogMasksCredentialShapedStderr`
  (`internal/agentctl/server/utility`, direct probe context)
- `TestStderrBufferTailKeepsTheEnd` and `TestStderrBufferBoundsRetentionOnWrite`
  (`internal/agentctl/server/utility`, bounded capture)
