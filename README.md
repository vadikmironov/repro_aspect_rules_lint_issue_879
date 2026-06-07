# repro_aspect_rules_lint_issue_879

<!-- --- BEGIN user-managed --- -->
Minimal reproduction for **aspect-build/rules_lint#879** — the Rust clippy
submodule split out in rules_lint 2.6.0 is not consumable by downstream repos.

In 2.5.2 clippy lived in the umbrella module
(`@aspect_rules_lint//lint:clippy.bzl`). 2.6.0 (#865) moved it into a separate
module, `aspect_rules_lint_rules_rust`, with the documented load
`@aspect_rules_lint_rules_rust//:clippy.bzl` — but that submodule cannot be
brought in by any downstream-consumable mechanism.

This repo isolates the three blockers as mutually-exclusive variants in
`MODULE.bazel`. **Variant A is active by default.** To run another, comment A and
uncomment B/C/D, then:

```bash
bazel mod graph        # forces module resolution; no build/toolchains needed
```

All variants are pinned at `main` HEAD `eff4e9396d` (2026-06-05) — i.e. *after*
dzbarsky's latest commit `30d965d` (PR #846) — to show the blockers persist.

| Variant | How it tries to pull in the submodule | Blocker | Expected error |
|---------|----------------------------------------|---------|----------------|
| **A** (active) | `bazel_dep` from BCR | #1 not published | `module aspect_rules_lint_rules_rust@0.0.0 not found in registries` |
| **B** | `archive_override` → GitHub `archive/<sha>.tar.gz` | #2 `export-ignore` strips `lint/rules_rust` from archives | `Prefix ".../lint/rules_rust" was given, but not found in the archive` |
| **C** | `git_override` (clone keeps the dir) | #3 submodule's versionless `bazel_dep` + dropped in-tree `local_path_override` | `bad bazel_dep on module 'aspect_rules_lint' with no version` |
| **D** | `archive_override` for the new `aspect_rules_lint_rust` module (PR #846) | #2 again — `.gitattributes` adds `lint/rust export-ignore` | `Prefix ".../lint/rust" was given, but not found in the archive` |

Notes:
- The root `aspect_rules_lint@2.6.0` resolves fine from BCR; only the Rust
  submodule is the problem.
- **B vs C** is the key contrast: same commit, same directory — `git_override`
  clones it successfully, yet `archive_override` can't find it. That isolates the
  cause to `export-ignore` (archive-only), not a missing feature.
- **Variant D** also requires switching the load in `tools/lint/linters.bzl` to
  `@aspect_rules_lint_rust//:clippy.bzl`, and that module depends on `rules_rs`
  (not `rules_rust`), so it is not a drop-in for a `rules_rust` repo.
<!-- --- END user-managed --- -->

## Project Layout

- `apps/` — your code: modules, services, and apps
- `MODULE.bazel` — module definition, segmented by language for easy pruning
- `tools/` — toolchains, formatters, and per-language dependency configuration

Add a module under `apps/`, then `bazel build //...` and
`bazel test //...` to confirm it wires up. The toolchains are already
configured — no per-language setup required.

## Install Bazelisk

Toolchains are hermetic, so Bazelisk is the only thing you need installed — it
reads `.bazelversion` and fetches the matching Bazel release on demand.

```bash
# Auto: system install (apt/.deb) when root or passwordless sudo is available,
# otherwise a no-sudo install into ~/.local/bin (PATH wired into your shell rc)
tools/setup/install_bazelisk.sh

# Force one mode or the other
tools/setup/install_bazelisk.sh --user      # no sudo, ~/.local/bin
tools/setup/install_bazelisk.sh --system    # apt/.deb, prompts for sudo
```

Linux only. Run `bazel version` to verify (restart your shell first if the
installer added `~/.local/bin` to your PATH).

## Common Commands

```bash
# Build / test everything
bazel build //...
bazel test //...                     # excludes lint tests — run lint separately

# List runnable targets — apps, buildifier, venv (plain `//...` lists everything)
bazel query 'kind("(py|cc|go|java|rust)_binary|buildifier|_venv", //...)'

# Format source (all languages), then Bazel/Starlark files
bazel run //:format
bazel run //:buildifier.fix          # Windows: tools\buildifier.bat fix
```

Linting is a separate, generated step — per-target `lint_test` rules are emitted
by Gazelle and then run as tests:

```bash
bazel run //:lint_gen                # preview without writing: -- -mode diff
bazel test --test_tag_filters=lint //...
```

### Regenerate dependency locks

After editing a language's dependency manifest, refresh its lockfile:

```bash
CARGO_BAZEL_REPIN=1 bazel fetch @crates//...                          # Rust   — tools/rust/Cargo.toml
bazel run @rules_go//go -- mod tidy                                   # Go     — go.mod / go.sum
```
