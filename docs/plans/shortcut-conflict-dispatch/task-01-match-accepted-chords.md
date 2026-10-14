---
id: "01-match-accepted-chords"
title: "Match accepted chords"
status: done
wave: 1
depends_on: []
plan: "plan.md"
requirements:
  - REQ-PLUGINS-PLUGIN-SHORTCUT-SETTINGS-001
acceptance_criteria:
  - AC-PLUGINS-PLUGIN-SHORTCUT-SETTINGS-001.5
system_design:
  - ../../specs/plugins/system-design/plugin-shortcut-settings.md
---

# Task 01: Match Accepted Chords

## Summary

Make existing conflict warnings and plugin core precedence agree with
`matchesShortcut` accepted events. Own pure overlap detection, event-specific
core reservation, and the existing conflict-label consumer in one sequential
implementation pass with independently authored regression evidence.

## In scope

- Replace `comboKey` and its callers with direct accepted-event overlap for
  warnings and actual-event core matching for dispatch.
- Keep direct two-entry groups; accumulate stable deduplicated peer labels in
  `buildConflictLabels` and reuse `useShortcutConflictLabels` on both cards.
- Preserve effective override reads, unbound behavior, fixed reservation
  defaults, and all unrelated dispatcher admission and handler behavior.
- Add focused tests to existing suites, with AC traceability annotations.

## Out of scope

Apply the plan's exclusions. Do not read, copy, replay, delete, or derive tests
from ROOT's private qualification source or `candidate.test.ts`/`source.json`.
No layout, copy, recorder, persistence, inventory, manifest, SDK, backend,
platform-matching, CodeMirror, or independent E2E changes are authorized.

## Acceptance

1. A fresh regression mounting real `usePluginShortcuts` and real
   `useKeyboardShortcut(SHORTCUTS.FIND_IN_PANEL)` fails before correction for
   the causal reason: Mac stub, plugin `ctrl+f`, cancelable ordinary Ctrl+F,
   core callback 0 and plugin callback 1. After correction core callback is 1,
   plugin callback 0, and default is prevented by core. Canonical Meta+F and
   noncolliding plugin Ctrl+Q controls pass before and after. Restore all
   platform stubs and unregister test handlers on cleanup.
2. Conflict and actual-label consumer tests show Mod overlaps Ctrl-only,
   Meta-only, and Ctrl+Meta but never make disjoint Ctrl-only and Meta-only
   actions label each other. Accumulate every direct peer once across groups.
   Cover mixed `ctrlOrCmd` with explicit `ctrl`/`cmd`, exact modifiers,
   key-case, different keys, Alt/Shift, and unbound entries using explicit
   behavioral expectations rather than tests mirroring the helper algorithm.
3. A core Ctrl-only override reserves Ctrl+Q while a Mod plugin still receives
   Meta+Q. Fresh override/removal, unbound, and fixed nonconfigurable core
   cases remain correct; a core-less reservation test observes no plugin
   `preventDefault`. Existing affected suites pass with no changes to
   matching compatibility, editor rules, registration order, or save flows.

## Implementation order

After a later explicit implementation release and one granted local-heavy
lease, mark this order `in_progress`. Ensure dependencies exist before pnpm
checks; only when missing, run `pnpm install --frozen-lockfile` from `apps/`.
Author the regression independently in the existing suite, run it RED with
controls, retain the causal failure receipt, then implement the minimum fix
and run GREEN. Add overlap, label-consumer, and partial-overlap regressions
before their production corrections. Run all product checks sequentially.
Record actual native handles, completion joins, exits, and current absence of
owned processes; retain failure evidence. Do not overlap heavy commands or
discard a native handle. Keep all work in this primary session.

## Verification

Run the RED selection before production edits, then repeat it for GREEN from
`apps/web`:

```bash
pnpm exec vitest run hooks/use-plugin-shortcuts.test.ts -t 'accepted-event core precedence' --maxWorkers=1
```

Name the new real-hook describe block `accepted-event core precedence` so this
selection includes the causal case and its two controls. Other RED selections
use the explicit newly added overlap/consumer case names in their owned suites;
record the exact commands and expected failures in Results.

Run these affected checks sequentially from `apps/web` after GREEN:

```bash
pnpm exec vitest run hooks/use-plugin-shortcuts.test.ts hooks/use-keyboard-shortcut.test.ts hooks/use-app-shortcuts.test.ts hooks/use-integration-shortcuts.test.ts lib/keyboard/shortcut-conflicts.test.ts lib/keyboard/utils.test.ts lib/keyboard/plugin-shortcuts.test.ts lib/keyboard/shortcut-overrides.test.ts components/settings/keyboard-shortcuts-card.test.tsx components/settings/plugins/plugin-shortcuts-card.test.tsx --maxWorkers=1
pnpm exec eslint --max-warnings 0 hooks/use-plugin-shortcuts.ts hooks/use-plugin-shortcuts.test.ts lib/keyboard/shortcut-conflicts.ts lib/keyboard/shortcut-conflicts.test.ts components/settings/keyboard-shortcuts-card.tsx components/settings/keyboard-shortcuts-card.test.tsx components/settings/plugins/plugin-shortcuts-card.test.tsx
pnpm run typecheck
```

After implementation is released, stage only task-owned changes for normal
delivery. Run `pnpm run i18n:ratchet` from `apps/web` with its normal active
selectors, including added/owned files through the index; no selector narrowing
or hook bypass. No new product copy or locale keys are planned. If actual owned
files differ, use their exact paths for ESLint and record that bounded change.

From repository root:

```bash
python3 scripts/list-docs.py validate
python3 scripts/lint-spec-files.py --all
git diff --check
git diff --cached --check
git status --short
```

Use this read-only working-tree coverage API preflight from repository root
with the activated Node toolchain. After implementation, replace the projected
source list only if actual owned paths changed; require `ok: true` and no errors.

```bash
node <<'NODE'
const fs = require('node:fs');
const { validateCoverage } = require('./.github/scripts/pr-docs.cjs');
const artifacts = [
  'docs/specs/plugins/requirements/plugin-shortcut-settings.md',
  'docs/specs/plugins/system-design/plugin-shortcut-settings.md',
  'docs/plans/shortcut-conflict-dispatch/plan.md',
  'docs/plans/shortcut-conflict-dispatch/task-01-match-accepted-chords.md',
];
const source = [
  'apps/web/lib/keyboard/shortcut-conflicts.ts',
  'apps/web/hooks/use-plugin-shortcuts.ts',
  'apps/web/components/settings/keyboard-shortcuts-card.tsx',
  'apps/web/lib/keyboard/shortcut-conflicts.test.ts',
  'apps/web/hooks/use-plugin-shortcuts.test.ts',
  'apps/web/components/settings/keyboard-shortcuts-card.test.tsx',
  'apps/web/components/settings/plugins/plugin-shortcuts-card.test.tsx',
];
const changedFiles = [...artifacts, ...source].map(filename => ({
  filename,
  status: filename.startsWith('docs/plans/shortcut-conflict-dispatch/') ? 'added' : 'modified',
}));
const fileContents = Object.fromEntries(artifacts.map(p => [p, fs.readFileSync(p, 'utf8')]));
const result = validateCoverage({ changedFiles, fileContents });
console.log(JSON.stringify(result, null, 2));
if (!result.ok || result.errors.length) process.exitCode = 1;
NODE
```

Do not invoke hosted status publication for a local dry-run. Record the command
and result. No broad local suite, build, QA, audit, or new mobile E2E is required
unless an actual failure makes a focused additional check necessary.

## Files likely touched

- `apps/web/lib/keyboard/shortcut-conflicts.ts`
- `apps/web/hooks/use-plugin-shortcuts.ts`
- `apps/web/components/settings/keyboard-shortcuts-card.tsx`
- `apps/web/lib/keyboard/shortcut-conflicts.test.ts`
- `apps/web/hooks/use-plugin-shortcuts.test.ts`
- `apps/web/components/settings/keyboard-shortcuts-card.test.tsx`
- `apps/web/components/settings/plugins/plugin-shortcuts-card.test.tsx`
- The four design-package artifacts, for accurate delivery status/results.

`utils.ts`, core reservation IDs, overrides, manifests, SDK, backend, and
plugin-detail production markup are read-only inputs.

## Dependencies

None. Later explicit implementation request and local-heavy lease required.

## Mobile and rendered scope

Reuse the current settings cards and translated warnings. This repair changes
shared data/action logic only. Attached keyboards on phones use the same
dispatch path. Real warning-consumer component tests provide the narrow
mobile-parity evidence; no phone layout preview or geometry E2E is warranted.
See the [plan](plan.md#mobile-and-rendered-scope).

## Risks

Nontransitive overlap, overwritten/duplicate labels, whole-binding suppression,
and platform stubs leaking between tests. Physical OS shortcut delivery remains
outside simulated keydown evidence.

## Parallelism

`sequential`

## Inputs

- Existing AC-PLUGINS-PLUGIN-SHORTCUT-SETTINGS-001.5 and its owning design's
  [conflict behavior](../../specs/plugins/system-design/plugin-shortcut-settings.md#conflict-and-dispatch-behavior).
- Repository matcher, real hook and registry, current core reservations, and
  affected settings-card tests. The qualification receipts establish the
  defect, not a test-source template.
- Plan exclusions and docs-maintainer assessment. Existing public precedence
  wording and completed relocation results remain valid.

## Results

- Frozen installation from `apps/` was required because this worktree had no
  dependencies: `pnpm install --frozen-lockfile` passed with pnpm 9.15.9.
- Independent RED: the exact `accepted-event core precedence` selection above
  exited 1 with one causal failure (core 0/plugin 1/default prevented) and two
  passing controls. GREEN repeated that command and passed all three tests.
- Additional RED:
  `pnpm exec vitest run hooks/use-plugin-shortcuts.test.ts lib/keyboard/shortcut-conflicts.test.ts components/settings/keyboard-shortcuts-card.test.tsx components/settings/plugins/plugin-shortcuts-card.test.tsx -t 'accepted overlap|direct conflict labels|event-specific core reservations' --maxWorkers=1`
  exited 1 with 16 expected failures and nine passes. This directly exposed
  accepted-event, label accumulation, rendered warning, and reservation defects.
- The exact ten-file affected Vitest command above passed 178/178 tests.
- The exact seven-file ESLint command initially exited 1 for an unused `isMac`
  import, a six-parameter test callback, and a 104-line test describe callback.
  Remove the import, retain the tuple through a rest callback, and split the
  describe block without changing behavior or assertions. The repeated scoped
  ESLint command passed with zero warnings.
- After those mechanical lint corrections:
  `pnpm exec vitest run hooks/use-plugin-shortcuts.test.ts lib/keyboard/shortcut-conflicts.test.ts components/settings/keyboard-shortcuts-card.test.tsx --maxWorkers=1`
  passed 59/59 tests. The other seven previously passing suites were unchanged.
- `NODE_OPTIONS=--max-old-space-size=8192 pnpm run typecheck` passed once before
  those mechanical lint corrections, including its normal generation and
  rich-output checks. No generated files changed.
- Normal `pnpm run i18n:ratchet` ran after staging all 11 owned files and passed
  (zero added and three modified production candidates; 643-entry allowlist).
  Test exclusions and all selectors remained unchanged.
- `python3 scripts/list-docs.py validate` passed (382 decisions, 1609 specs);
  `python3 scripts/lint-spec-files.py --all` passed; both working/index
  `git diff --check` commands passed. The exact coverage API heredoc above
  returned `ok: true`, `errors: []` for the actual 11 owned paths.
- Existing callback ordering, editor rules, fixed reservations, override
  behavior, matcher compatibility, save flows, and phone composition remain
  covered by the affected suites and unchanged production boundaries.
  Evidence uses simulated DOM keydown events; physical OS keys were not proven.
- All local native commands were joined and retained as external operational
  receipts with actual exits and current owned-group absence. Failed RED and
  lint receipts were kept. No private qualification source was read or replayed.
