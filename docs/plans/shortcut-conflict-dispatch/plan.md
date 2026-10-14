---
created: 2026-10-14
status: complete
requirements:
  - REQ-PLUGINS-PLUGIN-SHORTCUT-SETTINGS-001
system_design:
  - ../../specs/plugins/system-design/plugin-shortcut-settings.md
legacy_specs: []
---

# Implementation Plan: Shortcut Conflict Dispatch

## Overview

Restore AC-PLUGINS-PLUGIN-SHORTCUT-SETTINGS-001.5 by aligning conflict warnings
and plugin core precedence with the events `matchesShortcut` already accepts.
One sequential work order owns the matching, dispatch, and warning-consumer
changes together. Implementation followed the reviewed design checkpoint and
a later explicit request in the same primary session.

## Confirmed cause and evidence

`matchesShortcut` accepts Ctrl or Meta, including both, for `ctrlOrCmd` on
every platform. `comboKey` instead picks one platform modifier. The plugin
capture listener can therefore admit a plugin Ctrl+F handler on Mac, prevent
the event, and block the later real core Find bubble listener.

The supplied qualification receipts in
`/tmp/kandev-root-shortcut-precedence-20261014-rbnmzl_w/` report one causal
failure: plugin callback 1, core Find 0, `defaultPrevented=true`. Canonical
Meta+F and noncolliding Ctrl+Q controls passed; three selected tests took
23 ms. Native session 26924 joined at chunk `9acc65`, wrapper exit 0,
Vitest exit 1; group 2164896 was absent at the recorded check. The unsupported
`--minWorkers` attempt was setup failure, not causal evidence.

`post-4490-main-check.json` records all six implicated blobs equal between
tested `73c5ddb8c9550744655432d03edbe5872b6bbb80` and then-current main
`141c5a4743846abb090b40fbf2dd36cb15aae0a3`, with no test replay. These are
supplied simulated-event receipts; physical OS keys were not proven. The
private qualification source remains ROOT-owned and is not an implementation
input. Independently author permanent regression tests from repository
patterns and the observable contract.

## Scope

### In scope

- Existing accepted-event overlaps for effective shortcuts, exact modifiers,
  unbound entries, case-insensitive keys, and mixed modifier inputs.
- Event-specific precedence over matching plugin callbacks using current
  configurable core overrides and fixed core reservation defaults.
- Direct peer conflict labels accumulated without duplicate or transitive
  peers in the existing shared label consumer.
- Targeted pure, real-hook dispatch, and settings warning-consumer evidence.

### Out of scope

- Changing `matchesShortcut`, platform detection, persisted overrides,
  shortcut recording, manifest grammar, SDK, backend, or CodeMirror bindings.
- Changing eligibility, editor opt-in, IME, repeat/composition handling,
  registration order, multiple-plugin callbacks, or error handling.
- Expanding shortcut inventory, routes, layout, drafts, save flows, or
  subscriptions; manufactured mobile geometry checks or new E2E scenarios.
- Broad QA, security review, audits, full local suites, dependency updates,
  automatic delegation, or changes to completed relocation work orders.

## Technical approach

Follow [Conflict and dispatch behavior](../../specs/plugins/system-design/plugin-shortcut-settings.md#conflict-and-dispatch-behavior).
`shortcut-conflicts.ts` replaces platform canonicalization with bounded
accepted-event intersection and direct two-entry conflict groups.
`keyboard-shortcuts-card.tsx` accumulates and deduplicates true peer labels.
`use-plugin-shortcuts.ts` tests the actual event against effective core
bindings rather than comparing canonical plugin/core strings. Preserve fixed
reservations independently of overrides and preserve plugin events outside a
partial overlap. No public API or persistence change is required.

## Tests

All new cases trace to AC-PLUGINS-PLUGIN-SHORTCUT-SETTINGS-001.5:

- `hooks/use-plugin-shortcuts.test.ts`: independently authored real-hook
  regression `preserves explicit Ctrl+F on Mac`, canonical
  Meta+F and free Ctrl+Q controls, partial-overlap dispatch, fresh overrides,
  unbound core/plugin behavior, and fixed reservation defaults.
- `lib/keyboard/shortcut-conflicts.test.ts`: direct overlap/nonoverlap cases
  for Ctrl, Meta, both, Mod, mixed explicit fields, Alt/Shift, keys/case, and
  unbound entries. Update obsolete platform-canonicalization expectations.
- `components/settings/keyboard-shortcuts-card.test.tsx`: actual label
  consumer and rendered core warning for Ctrl/Meta/Mod bridges; accumulate
  peers exactly once and exclude disjoint peers.
- `components/settings/plugins/plugin-shortcuts-card.test.tsx`: rendered
  selected-plugin warning uses the same effective overlap and qualified
  labels, including direct-only peer names.
- Existing matcher, override, plugin-entry, app/integration shortcut, and
  keyboard-hook suites preserve unaffected matching and admission behavior.

## Mobile and rendered scope

This is state/data and action dispatch within existing surfaces. Both settings
cards retain their current translated warning template and composition.
Phone attached-keyboard events use the same hooks and conflict lookup.
Targeted component tests exercise actual rendered warning consumers; no new
geometry, navigation, or touch contract calls for an ASCII redesign or mobile
Playwright scenario under the narrow mobile-parity exception.

The completed [relocation plan](../plugin-shortcut-settings/plan.md) remains
historical evidence for member access, save/reload, and phone detail entry in
`e2e/tests/auth/plugin-settings-member.spec.ts` and
`e2e/tests/plugins/mobile-plugin-settings-row.spec.ts`. Do not rerun or change
those scenarios solely for this data/dispatch repair.

## Public documentation impact

Docs-maintainer assessment: `docs/public/plugins.md` already says core
shortcuts always win; `docs/public/plugins-authoring.md` already documents the
same precedence and personal overrides. Root README and screenshot catalog
need no update. The repair restores that existing contract without new copy,
configuration, screenshots, or author API. Update only the owning internal
design and its delivery links; no public-doc polish is required.

## Work orders

- [x] [Task 01: Match accepted chords](task-01-match-accepted-chords.md)

## Verification results

The design checkpoint passed
`python3 scripts/list-docs.py validate` (382 decisions, 1609 specifications),
`python3 scripts/lint-spec-files.py --all`, scoped whitespace checks, and the
documentation-coverage API dry-run (`ok: true`, `errors: []`) with the four
artifacts and seven projected owned source/test paths. No hosted status was
published. Task 01 owns all later product checks and their exact results.

The shell did not initially expose Node on PATH; the coverage dry-run was
completed using the existing Node 24.18.0 binary, without installing anything.
Activate the repository's Node 24 / pnpm 9.15.9 toolchain before the later
work-order package commands. No heavy command was started at the design checkpoint.

Task 01 is complete. Independent dispatch RED reproduced one causal failure
with two passing controls; GREEN passed all three. Overlap/label/reservation
RED retained 16 expected failures and nine passes. The ten affected suites
passed all 178 tests. Scoped ESLint passed after an unused-import cleanup and
two mechanical test-structure corrections; the three changed suites then
passed all 59 tests. The single 8 GiB typecheck passed before those mechanical
lint corrections. Normal staged i18n ratchet, spec catalog, all-spec lint,
coverage API preflight, and working/index whitespace checks passed. Exact
commands and results are in the completed work order. No public-doc, mobile
layout/E2E, matcher, or persistence changes were needed.

## Risks

- Accepted-event overlap is not transitive; connected-component grouping would
  falsely label Ctrl-only and Meta-only actions through a Mod bridge.
- Multiple pair groups can overwrite labels or duplicate peer names unless
  the actual consumer accumulates them.
- Whole-binding suppression would incorrectly discard plugin events outside a
  partial core overlap. Platform-specific matching would break compatibility.
- Simulated events establish DOM routing only. Do not claim physical OS key
  delivery or reinterpret another process's receipt as a local test run.
