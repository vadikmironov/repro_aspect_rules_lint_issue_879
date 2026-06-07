"""
This module contains all linters definition as per the rules_lint documentation
https://github.com/aspect-build/rules_lint/blob/main/docs/linting.md
"""

# rules_lint 2.6.0 (#865) moved clippy out of the umbrella module into the
# aspect_rules_lint_rules_rust submodule, so the documented load path is now
# @aspect_rules_lint_rules_rust//:clippy.bzl (the pre-2.6.0 path
# @aspect_rules_lint//lint:clippy.bzl no longer exists). This is the load that
# can't be satisfied downstream — see MODULE.bazel and the issue #879 repro.
load("@aspect_rules_lint//lint:lint_test.bzl", "lint_test")
load("@aspect_rules_lint_rules_rust//:clippy.bzl", "lint_clippy_aspect")

clippy = lint_clippy_aspect(
    config = Label("@//:.clippy.toml"),
)

clippy_test = lint_test(aspect = clippy)
