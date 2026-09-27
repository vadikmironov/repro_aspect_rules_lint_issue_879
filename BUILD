load("@buildifier_prebuilt//:rules.bzl", "buildifier")
load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "nogo")

package(default_visibility = ["//visibility:public"])

# --- BEGIN user-managed ---
# Repo-specific Bazel/Gazelle customizations — preserved across re-bootstrap.
# Add gazelle:exclude directives and tweak the buildifier excludes here.
#
# NOTE: the buildifier_prebuilt macro joins exclude_patterns with `-o` and
# appends ONE trailing `-prune`; find binds that prune to only the last
# `-path`, so keep this to a single pattern if you need pruning to take effect.
_BUILDIFIER_EXCLUDES = ["./.git/*"]
# --- END user-managed ---

# Buildifier lints and autoformats bazel (Starlark) files.
#
#   bazel run //:buildifier.fix
#   bazel run //:buildifier.check
#
# Windows works from buildifier_prebuilt 8.5.1.4 on, which fixed the runner's
# argument escaping (keith/buildifier-prebuilt#168).

buildifier(
    name = "buildifier.check",
    exclude_patterns = _BUILDIFIER_EXCLUDES,
    lint_mode = "warn",
    mode = "diff",
)

buildifier(
    name = "buildifier.fix",
    exclude_patterns = _BUILDIFIER_EXCLUDES,
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
