# repro_aspect_rules_lint_issue_879

<!-- --- BEGIN user-managed --- -->
Tracks the state of downstream Rust (clippy) linting in **aspect_rules_lint**,
for a repo that builds its Rust with **BCR `rules_rust`**. Originally a repro for
[#879](https://github.com/aspect-build/rules_lint/issues/879); now it shows where
2.7.1 lands.

## The #879 packaging blockers are fixed in 2.7.1

In 2.6.0 the clippy aspect moved out of the umbrella module into a separate
module that downstream repos couldn't pull in (not on the BCR; `export-ignore`
stripped it from release archives; an in-tree `local_path_override` made
overrides unresolvable). As of **2.7.1** that module is published to the Bazel
Central Registry as **`aspect_rules_lint_rust`** (renamed from
`aspect_rules_lint_rules_rust`), so a plain `bazel_dep` resolves with no
overrides — this is the active **Variant A** in `MODULE.bazel`:

```starlark
bazel_dep(name = "aspect_rules_lint", version = "2.7.1")
bazel_dep(name = "aspect_rules_lint_rust", version = "0.0.2")
```

```bash
bazel mod graph        # resolves; pulls aspect_rules_lint@2.7.1 + rules_rs@0.0.83
```

(The old packaging-blocker variants B/C/D — `archive_override` / `git_override`
against `main` HEAD — are kept commented in `MODULE.bazel` as pre-2.7.1 history.)

## …but the clippy aspect silently no-ops on a BCR `rules_rust` target

`apps/rust_app/src/main.rs` contains a deny-by-default `clippy::eq_op` violation,
yet the lint test passes green:

```bash
bazel test //apps/rust_app:rust_app.lint --test_tag_filters=lint
# //apps/rust_app:rust_app.lint   PASSED      <-- expected FAIL
```

clippy never runs. The report artifact
(`bazel-bin/apps/rust_app/rust_app.AspectRulesLintClippy.report`) is empty, and
`bazel test … -s` shows the aspect's action is just
`touch …report && echo 0 > …exit_code`. The toolchain's own `clippy-driver` does
flag the code, so the toolchain is fine:

```
error: equal expressions as operands to `==`   (#[deny(clippy::eq_op)])
```

### Why

`aspect_rules_lint_rust`'s `MODULE.bazel` sources `@rules_rust` from `rules_rs`
(`use_extension("@rules_rs//rs:rules_rust.bzl", …)`), so the aspect's
`rust_clippy_action.get_clippy_ready_crate_info(target, ctx)` (`clippy.bzl:134`)
looks for a `CrateInfo` from `rules_rs++rules_rust+rules_rust`. This repo builds
`rust_app` with **BCR `rules_rust@0.70.0`** (`rules_rust+`), whose `CrateInfo` is
a different provider. `crate_info` comes back `None`, so the aspect takes its
`noop_lint_action` branch (`clippy.bzl:144`). Both `rules_rust` repos coexist in
the module graph.

Net: 2.7.1 is consumable, but only actually lints `rules_rs`-built targets — for
a BCR `rules_rust` repo it's a false-green. Write-up in `tmp/issue_879_comment.md`.
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
