---
id: "16-complete-ci-delivery"
title: "Measure and deliver complete CI"
status: in_progress
wave: 12
depends_on: ["15-remaining-critical-paths"]
plan: "plan.md"
requirements:
  - REQ-PLATFORM-CI-PERFORMANCE-003
  - REQ-PLATFORM-CI-PERFORMANCE-007
acceptance_criteria:
  - AC-PLATFORM-CI-PERFORMANCE-007.1
  - AC-PLATFORM-CI-PERFORMANCE-007.2
  - AC-PLATFORM-CI-PERFORMANCE-007.3
  - AC-PLATFORM-CI-PERFORMANCE-007.4
system_design:
  - ../../specs/platform/system-design/ci-completion-budget.md
---

# Task 16: Measure and deliver complete CI

## Outcome and scope

Collect hosted attempts, exact coverage reports, workflow identities and occupancy. Keep one bounded rollout active; resolve valid CI/review failures and merge the final passing head.

## Exclusions

No new product behavior, paid capacity activation, trust expansion, assertion
removal, automatic retries or delegation. Repair a demonstrated existing-contract
regression only when needed to retain a failing CI assertion. Preserve unrelated user changes.

## Acceptance

Three complete successful workflow sets meet 900 seconds including waiting. Partial, cancelled and failed runs remain visible and do not count. All required checks and reviews pass at the final head before authorized squash merge.

## Files and dependencies

The backend/frontend/E2E workflow files, their registered contract tests and
CI performance evidence. Task 15 supplies the preceding measured result.

## Verification

```bash
python3 .github/scripts/ci-efficiency_test.py
scripts/pr-state --compact <PR>
scripts/pr-resolve list <PR>
python3 scripts/list-docs.py validate
python3 scripts/lint-spec-files.py --all
git diff --check
```

## Results

Hosted evidence remains mandatory; local contracts do not prove speed.

Second candidate head `e2a66b8536` runs backend 38581852307, frontend
38581852072 and E2E 38581852055. Its tested merge is
`048610fb914664e1873c4ac671c66c1622bde238` over base `73c5ddb8c955`.
The initial head `6025705717` failed frontend argument forwarding and two cold
Windows compilation steps; the subsequent push superseded its remaining jobs.
Record that attempt as failed/partial, never complete timing acceptance.

For controlled frontend comparisons, a temporary branch
`perf/full-ci-frontend-baseline` starts at the exact tested merge. Commit
`15b7acfc46381193ab6f792a43ad2b9de6f0d3bd` restores the base frontend workflow
and adds only a JSON report, required artifact upload and workflow/checkout
identity recording. Both candidate and baseline application tree hashes are
`c3cfe17ff326562b01df5257bedbdf92fc23c511`. Its normal hooks passed with no
bypass. Dispatch 38582284835 is the first unsharded measurement. This branch
has no PR and will not merge; comparisons must retain event/cache differences.

Review thread 4271878983 correctly identified an overly broad frontend job
range in the runner placement contract. Narrow it to the static job and
explicitly check the new repository job's hosted boundary. Five placement
contract tests pass; thread disposition follows the verified fix commit.

The first complete frontend candidate succeeded in 825 seconds (13m45s), with
757 seconds of execution critical path and 30.15 runner-minutes. Its two reports
contain 26,411 tests (26,407 passed, four skipped) across all 2,948 files. The
normal repository artifact union is exactly 1,904 top-level tests, once each,
at merge source `048610fb914664e1873c4ac671c66c1622bde238`.

A matching temporary branch `perf/full-ci-frontend-candidate` points directly to
that tested merge. Dispatch-only frontend comparisons on the two temporary
branches preserve identical application trees and avoid PR-event differences.
The first PR result is useful loaded evidence; baseline reports and repeated
matched dispatch samples remain pending.

Third candidate `5949b632b0` passed the complete backend workflow 38587972467
in 1,134 seconds (18m54s), with 971 seconds of execution critical path and
243.32 runner-minutes. All four Windows process reports validate the same native
inventory: 864 names across two packages, split 216/216/216/216 exactly once.
Frontend 38587972660 failed only the existing Kubernetes SemVer workflow-string
contract, which still expected the old make target; update that assertion to the
new versioned native build command. E2E 38587972415 remains pending and cannot
count as complete acceptance. Its build improved from approximately eight to
six minutes.

Matched frontend dispatch pair two (38587875054 / 38587880385) passed with
identical 26,411 test identities and outcomes. Full elapsed was 1,725 versus
704 seconds; execution critical path was 1,540 versus 663 seconds; runner-minutes
were 25.6667 versus 26.5333 (+3.38%). Pair three is running as
38592146394 / 38592152123 on the same frozen application tree. The earlier PR
sample's higher occupancy remains retained and is not hidden by these pairs.

The fourth candidate's local verification covers canonical shared-manifest
identity/digest checks, isolated caches and failed-preparation cleanup, actual
worker-image command boundaries and cleanup, existing Kubernetes pins, native
Go allocation, cache/workflow contracts and unchanged process resource guards.
The real local Kubernetes worker-preset lifecycle passed for minimal, node-pnpm
and Python images (one test, 5.2 minutes including fixture setup; test body 37.7
seconds). With the real 44 MiB helper, preparation measured 21.894 seconds once
and 0.643/0.626/0.664 seconds for three isolated bundles. These local measurements
are diagnostic; the complete hosted acceptance gate remains unmet.

The same 3,808-test normal manifest partitions into 84 disjoint one-worker
processes across 28 jobs. Its longest whole-file estimate is 438.6825 seconds,
so the idealized 343-second average is not the expected tail. Native ordinary Go
selection remains exactly 407 packages in buckets 97/105/112/93, with the
repository package delegated once to its required native-cohort matrix.

### Third-candidate mobile failures and bounded remediation

E2E run 38587972415 failed in 2,076 seconds (34m36s), with 1,767 seconds
execution critical path and 377.43 runner-minutes. Preserve this failed attempt.
Its two failures are mobile empty-turn hints and non-active parent defaults.
The trace's sidebar response contains the selected parent's repository, but the
subtask dialog reads only the active Kanban collection. A browser connected
before the parent finishes can mask that missing data dependency through live
events. The dialog's initial repository effect then runs only on mount.

Keep the existing mobile composition and assertions. Pass the selected drawer
row into the existing subtask dialog as initial parent context; prefer a matching
current store row when available. Do not change persisted settings or form edits.
Make the existing E2E wait for parent completion before navigation to prove the
paged-row path without live creation events. Cover explicit selected-parent
context independently of the active task and preserve desktop callers.

The empty-turn notice is live-only. Use the same established desktop test
sequence on mobile: open a completed ordinary turn, then send the empty slash
command after chat readiness. Do not lengthen retries/timeouts or remove the
hint assertion. Verify both rendered mobile scenarios together after rebuilding.
This is CI regression remediation under the user's continuous delivery authority.

The strengthened selected-parent test failed locally before the fix, matching the
hosted missing-chip assertion. After the row handoff and a fresh production web
build, all three focused mobile scenarios passed in 20.3 seconds: the live hint
(using the phone Send button), ordinary Create Subtask, and completed non-active
parent defaults. Eleven existing component/context tests, changed-file ESLint,
formatting, full web typecheck and documentation validation passed. No UI layout,
copy, public API, or operator instructions change; this restores existing defaults.

Desktop Create Subtask also passed (one test, ten seconds including setup).
The normal commit hook enforced the component-length limit; extracting the
mobile target state/handler into a focused hook preserves the same row handoff.
The extraction passes strict lint and the seven mobile component tests.
Matched frontend candidate three passed all 26,411 tests in 1,074 seconds full
elapsed, 739 seconds execution critical path and 30.48 runner-minutes. Its paired
unsharded baseline remains running. The added waiting is retained in the result.
