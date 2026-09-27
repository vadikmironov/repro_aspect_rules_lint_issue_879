# repro_aspect_rules_lint_issue_879

<!-- --- BEGIN user-managed --- -->
Minimal repro for
[aspect-build/rules_lint#879](https://github.com/aspect-build/rules_lint/issues/879):
the clippy aspect from `aspect_rules_lint_rust` passes a Rust target that is
built with BCR `rules_rust`, and clippy does not run.

Last checked on 2026-09-27 with Bazel 9.2.0, `aspect_rules_lint` 2.9.0,
`aspect_rules_lint_rust` 0.0.3 and `rules_rust` 0.74.0.

## Packaging is fixed

The packaging problems in the issue title are fixed since 2.7.1.
`aspect_rules_lint_rust` is on the BCR, and a plain `bazel_dep` resolves:

```starlark
bazel_dep(name = "aspect_rules_lint", version = "2.9.0")
bazel_dep(name = "aspect_rules_lint_rust", version = "0.0.3")
```

## The clippy aspect does not run on a BCR `rules_rust` target

`apps/rust_app/src/main.rs` has a deny-by-default `clippy::eq_op` violation.
The lint test passes:

```bash
bazel test //apps/rust_app:rust_app.lint --test_tag_filters=lint
# //apps/rust_app:rust_app.lint   PASSED      <-- expected FAIL
```

The report `bazel-bin/apps/rust_app/rust_app.AspectRulesLintClippy.report` is
empty. The aspect action only creates the report and writes exit code 0:

```bash
bazel aquery 'outputs(".*AspectRulesLintClippy.report", deps(//apps/rust_app:rust_app))' \
    --aspects=//tools/lint:linters.bzl%clippy --output_groups=rules_lint_human
# touch …/rust_app.AspectRulesLintClippy.report && echo 0 > …/rust_app.AspectRulesLintClippy.report.exit_code
```

The toolchain is correct. The clippy aspect of `rules_rust` finds the violation:

```bash
bazel build //apps/rust_app:rust_app \
    --aspects=@rules_rust//rust:defs.bzl%rust_clippy_aspect --output_groups=clippy_checks
# error: equal expressions as operands to `==`   (`#[deny(clippy::eq_op)]` on by default)
```

### Why

`aspect_rules_lint_rust` gets `@rules_rust` from `rules_rs`
([MODULE.bazel](https://github.com/aspect-build/rules_lint/blob/rust-v0.0.3/lint/rust/MODULE.bazel#L25)).
This repo builds `rust_app` with BCR `rules_rust`. The two modules define
different `CrateInfo` providers, so the aspect finds no `CrateInfo` on
`rust_app`
([clippy.bzl#L134](https://github.com/aspect-build/rules_lint/blob/rust-v0.0.3/lint/rust/clippy.bzl#L134))
and uses its no-op action
([clippy.bzl#L144-L145](https://github.com/aspect-build/rules_lint/blob/rust-v0.0.3/lint/rust/clippy.bzl#L144-L145)).

## Test a local rules_lint clone

Override both modules with a git clone of rules_lint. Archives made by
`git archive` do not contain `lint/rust`.

```bash
bazel test //apps/rust_app:rust_app.lint --test_tag_filters=lint --lockfile_mode=off \
    --override_module=aspect_rules_lint=$HOME/src/rules_lint \
    --override_module=aspect_rules_lint_rust=$HOME/src/rules_lint/lint/rust
```
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

Downloads go through `tools/setup/download_lib.sh`, which resumes from the bytes
already on disk instead of restarting when a transfer stops making progress.
Behind a slow or stalling proxy, raise `DOWNLOAD_STALL_SECONDS` (default 30, the
seconds under `DOWNLOAD_MIN_SPEED` before an attempt is abandoned) or
`DOWNLOAD_MAX_STALLED_RETRIES` (default 5).

## Common Commands

```bash
# List runnable targets — apps, buildifier, venv (plain `//...` lists everything)
bazel query 'kind("(py|cc|go|java|rust)_binary|buildifier|_venv", //...)'

# Build / test everything
bazel build //...
bazel test //...                     # excludes lint tests — run lint separately

# Format source (all languages), then Bazel/Starlark files
bazel run //:buildifier.fix
bazel run //:format
```

Linting is a separate, generated step — per-target `lint_test` rules are emitted
by Gazelle and then run as tests:

```bash
bazel run //:lint_gen                # preview without writing: -- -mode diff
bazel test --test_tag_filters=lint //...
```

### Sync dependencies

After adding or removing a module extension repo in a `*.MODULE.bazel` segment, sync
the `use_repo()` calls.
After editing a language's dependency manifest, refresh its lockfile too.

```bash
bazel mod tidy                                                        # Bazel  — use_repo() calls in the MODULE segments
CARGO_BAZEL_REPIN=1 bazel fetch @crates//...                          # Rust   — tools/rust/Cargo.toml
bazel run @rules_go//go -- mod tidy                                   # Go     — go.mod / go.sum
```

## Local Disk Cache

Faster local builds, shared across every Bazel project you build. Add to your
**user-global `~/.bazelrc`** (not this repo — it applies to all your workspaces):

```bash
build --disk_cache=~/.cache/bazel-disk
build --experimental_disk_cache_gc_max_size=15G   # bounded; auto-GC'd when idle (Bazel 7.4+)
```

The same action-cache (AC/CAS) mechanism as the remote cache, on local disk:
identical actions (e.g. a shared protobuf compile) run once and are reused
everywhere, and `bazel clean` becomes cheap to recover from. Trades disk for
speed — tune the size to taste.
