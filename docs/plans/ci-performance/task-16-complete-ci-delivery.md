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

No application behavior, paid capacity activation, trust expansion, assertion
removal, automatic retries or delegation. Preserve unrelated user changes.

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
