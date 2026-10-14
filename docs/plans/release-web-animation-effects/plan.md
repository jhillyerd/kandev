---
created: 2026-10-14
status: done
requirements:
  - REQ-UI-CHAT-MOTION-003
system_design:
  - ../../specs/ui/system-design/chat-motion.md
legacy_specs: []
---

# Implementation Plan: Release Web Animation Effects

## Overview

Stop the macOS desktop web view from retaining every animated chat row and text
run. Each script animation owner detaches the effect when it cancels the animation.
One work order adds the shared helper, applies it at every `animate()` owner, and
adds regression assertions.

This package follows the heap measurements in
[issue 4100](https://github.com/kdlbs/kandev/issues/4100#issuecomment-6153587005).

## Requirement conformance

AC-UI-CHAT-MOTION-003.4 already requires that sustained streaming does not
create an ever-growing backlog of effects. No new requirement is necessary. The
[system design](../../specs/ui/system-design/chat-motion.md) now states the
release rule. This creates no new ownership, persistence, or API boundary, so no
ADR is needed.

## Evidence and confirmed root cause

On desktop 0.97.0 (macOS 27.0), `heap` on the WebContent process showed
`KeyframeEffect`, `WebAnimation`, `HTMLSpanElement`, and `HTMLDivElement` growing
together and never dropping: 1,335 effects at 17 minutes, 13,527 at 4 hours, and
25,438 at 23 hours. `CSSAnimation` stayed under 35 and `CSSStyleSheet` at 10.

A standalone WKWebView harness with no Kandev code reproduced it. Each case
appended 2,000 spans, then applied GC pressure:

| Case | Effects retained | Spans retained |
| --- | --- | --- |
| Removed without animation | 0 | 0 |
| CSS animation, then removed | 0 | 0 |
| `animate()`, `cancel()`, removed | 2,000 | 2,000 |
| `animate()` finished, removed | 2,000 | 2,000 |
| `cancel()`, `effect = null`, removed | 0 | 0 |

Only the small `WebAnimation` objects remain in the last case. Headless Chromium
against the same app did not retain animations.

`ChatMotionSpan` animates each appended text run, and `ChatMotionItem` animates
each new row. Both canceled correctly, but WebKit still kept the detached
elements and their subtrees.

## Scope

### In scope

- `releaseAnimation(animation)` in `apps/packages/ui/src/animation-utils.tsx`.
- Use it in chat row and text motion, grid spinner, compositor spin and pulse,
  and settings-search motion.
- Unit assertions that chat motion detaches the effect on cleanup.

### Out of scope

- Switching chat motion to CSS animations.
- The upstream WebKit defect.
- Store retention, stylesheet, and scroll-loop changes.

## Work orders

- [x] [Task 01: Release animation effects](task-01-release-animation-effects.md)

## Verification results

- The two new assertions failed with only `cancel()` and pass with `releaseAnimation`.
- The chat motion, grid spinner, compositor, motion visibility, and chat-motion
  hook suites pass (12 files, 123 tests). Web typecheck and scoped ESLint pass.
- The WKWebView harness released all effects and spans with the same cleanup.

## Risks

- Setting `effect` to null on an already canceled animation has no visible effect.
- `jsdom` does not implement Web Animations, so tests use handle doubles.

## Documentation impact

Internal design and delivery documents change. Public documentation does not
change because behavior, settings, and controls stay the same.
