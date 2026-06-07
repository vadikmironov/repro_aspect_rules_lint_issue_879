load("@buildifier_prebuilt//:rules.bzl", "buildifier")
load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "nogo")

package(default_visibility = ["//visibility:public"])

# Buildifier lints and autoformats bazel (Starlark) files.
#
# Mac/Linux: Use the bazel targets directly
#   bazel run //:buildifier.fix
#   bazel run //:buildifier.check
#
# Windows: The macro has compatibility issues, use the wrapper script instead
#   tools\buildifier.bat fix
#   tools\buildifier.bat check
#
# See [tools/buildifier.md](tools/buildifier.md) for details on the Windows workaround.

buildifier(
    name = "buildifier.check",
    exclude_patterns = ["./.git/*"],
    lint_mode = "warn",
    mode = "diff",
)

buildifier(
    name = "buildifier.fix",
    exclude_patterns = ["./.git/*"],
    lint_mode = "fix",
    mode = "fix",
)

# rules_lint integration - export linter configurations
alias(
    name = "format",
    actual = "//tools/format",
)

exports_files(
    [
        ".rustfmt.toml",
        ".clippy.toml",
    ],
    visibility = ["//visibility:public"],
)

# rules_go configuration

# Go module definition exported for use in go_segment.MODULE.bazel.
# Reference: https://go.dev/doc/modules/gomod-ref
exports_files(["go.mod"])

# Static analysis tool for Go code (nogo; part of the lint feature). Registered
# via go_sdk.nogo in tools/go/go_segment.MODULE.bazel.
# Reference: https://github.com/bazel-contrib/rules_go/blob/master/go/nogo.rst
nogo(
    name = "repro_aspect_rules_lint_issue_879_nogo",
    config = ":.nogo_config.json",
    vet = True,
    visibility = ["//visibility:public"],
)

# Gazelle driver for the lint extension. Auto-generates lint_test
# targets (ruff_test for py_*, etc.) tagged "lint" so CI can select
# them via `bazel test --test_tag_filters=lint //...`. See
# tools/lint/gazelle/.
#   bazel run //:lint_gen                  # apply
#   bazel run //:lint_gen -- -mode diff    # preview
gazelle(
    name = "lint_gen",
    gazelle = "//tools/lint/gazelle:gazelle_lint",
)
