import assert from "node:assert/strict";
import { test } from "node:test";
import { cacheIdentity } from "./go-cache-identity.mjs";

const identity = {
  family: "test-1",
  day: "2026-10-14",
  runner: "Linux-X64",
  dependencies: "dependency-hash",
  recipe: "image-recipe-hash",
  go: {
    GOVERSION: "go1.26.0",
    GOOS: "linux",
    GOARCH: "amd64",
    CGO_ENABLED: "1",
    GOFLAGS: "-trimpath",
    GOAMD64: "v1",
    GOEXPERIMENT: "",
    GOMODCACHE: "/root/go/pkg/mod",
    GOCACHE: "/cache with spaces/build",
  },
  compiler: "gcc 13.3.0",
};
// @covers AC-PLATFORM-CI-PERFORMANCE-007.1, 007.2, 007.3

test("daily rollover preserves compatible restore prefix", () => {
  const today = cacheIdentity(identity);
  const tomorrow = cacheIdentity({ ...identity, day: "2026-10-15" });
  assert.notEqual(today.buildKey, tomorrow.buildKey);
  assert.equal(today.buildPrefix, tomorrow.buildPrefix);
  assert.ok(today.buildKey.startsWith(tomorrow.buildPrefix));
  assert.equal(today.moduleKey, tomorrow.moduleKey);
});

test("same day reuses the generation", () => {
  assert.deepEqual(
    cacheIdentity(identity),
    cacheIdentity({ ...identity, commit: "source-only-change" }),
  );
  assert.equal(cacheIdentity(identity).buildPath, "/cache with spaces/build");
});

test("compatibility changes cannot match", () => {
  const original = cacheIdentity(identity);
  for (const change of [
    { family: "test-2" },
    { runner: "Linux-ARM64" },
    { dependencies: "new-deps" },
    { recipe: "new-image" },
    { compiler: "gcc 14" },
    ...[
      "GOVERSION",
      "GOOS",
      "GOARCH",
      "CGO_ENABLED",
      "GOFLAGS",
      "GOAMD64",
      "GOEXPERIMENT",
    ].map((key) => ({ go: { ...identity.go, [key]: "different" } })),
  ]) {
    const changed = cacheIdentity({ ...identity, ...change });
    assert.ok(
      !original.buildKey.startsWith(changed.buildPrefix),
      JSON.stringify(change),
    );
  }
  assert.equal(
    original.moduleKey,
    cacheIdentity({ ...identity, family: "static" }).moduleKey,
  );
});

test("invalid identity disables optional caching", () => {
  for (const change of [
    { dependencies: "" },
    { recipe: "" },
    { family: "../invalid" },
    { day: "2026-02-30" },
    { go: { ...identity.go, GOCACHE: "off" } },
    { go: { ...identity.go, GOMODCACHE: "/tmp/bad\npath" } },
    { go: { ...identity.go, GOCACHE: "/tmp/*" } },
    { compiler: "" },
  ])
    assert.throws(() => cacheIdentity({ ...identity, ...change }), /invalid/i);
});

test("linter version and configuration isolate lint snapshots", () => {
  const lint = {
    version: "golangci-lint 2.9.0",
    config: "lint-config",
    path: "/root/.cache/golangci-lint",
  };
  const current = cacheIdentity({ ...identity, family: "static", lint });
  for (const changed of [
    { ...lint, version: "2.10.0" },
    { ...lint, config: "new-config" },
  ]) {
    assert.ok(
      !current.lintKey.startsWith(
        cacheIdentity({ ...identity, family: "static", lint: changed })
          .lintPrefix,
      ),
    );
  }
});

test("CLI identity failure leaves required commands runnable", async () => {
  const { mkdtempSync, readFileSync, rmSync } = await import("node:fs");
  const { tmpdir } = await import("node:os");
  const { spawnSync } = await import("node:child_process");
  const directory = mkdtempSync(`${tmpdir()}/go-cache-`);
  try {
    const result = spawnSync(
      "bash",
      ["-e", "-c", 'node "$IDENTITY_SCRIPT"; echo required-command-ran'],
      {
        encoding: "utf8",
        env: {
          ...process.env,
          IDENTITY_SCRIPT: new URL("./go-cache-identity.mjs", import.meta.url)
            .pathname,
          GITHUB_WORKSPACE: directory,
          GITHUB_OUTPUT: `${directory}/output`,
          CACHE_FAMILY: "test-1",
        },
      },
    );
    assert.equal(result.status, 0, result.stderr);
    assert.match(result.stdout, /required-command-ran/);
    assert.equal(
      readFileSync(`${directory}/output`, "utf8"),
      "enabled=false\n",
    );
  } finally {
    rmSync(directory, { recursive: true, force: true });
  }
});
