"""
This module contains all linters definition as per the rules_lint documentation
https://github.com/aspect-build/rules_lint/blob/main/docs/linting.md
"""

# clippy ships in the standalone aspect_rules_lint_rust module (BCR), loaded from
# @aspect_rules_lint_rust//:clippy.bzl.
load("@aspect_rules_lint//lint:lint_test.bzl", "lint_test")
load("@aspect_rules_lint_rust//:clippy.bzl", "lint_clippy_aspect")

clippy = lint_clippy_aspect(
    config = Label("@//:.clippy.toml"),
)

clippy_test = lint_test(aspect = clippy)
