# gazelle:lint_ignore
#
# This BUILD defines the linter-tool binaries themselves (clang_tidy,
# pmd, spotbugs) as wrapper rules around prebuilt JARs/binaries. They
# carry no source files, so spotbugs/pmd/clang_tidy aspects produce
# no rules_lint_human output group, and any auto-generated lint_test
# would fail analysis. lint_ignore opts the entire package out.

#
# This BUILD file is derived from the aspect_rules_lint example:
# https://github.com/aspect-build/rules_lint/blob/main/example/tools/lint/BUILD.bazel
#

package(default_visibility = ["//:__subpackages__"])
