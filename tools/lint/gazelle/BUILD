load("@gazelle//:def.bzl", "gazelle_binary")
load("@rules_go//go:def.bzl", "go_library")

# Gazelle language extension that emits lint_test targets. See lang.go
# for the extension's contract and behaviour. Per-language generators
# in cpp.go, rust.go, java.go, python.go are wrapped in lang-segment
# markers so a scaffolded fork that omits a language compiles this
# extension with that language's generator absent.
go_library(
    name = "lint_lang",
    srcs = [
        "generate.go",
        "kinds.go",
        "lang.go",
        "rust.go",
    ],
    importpath = "repro_aspect_rules_lint_issue_879/tools/lint/gazelle",
    visibility = ["//visibility:public"],
    deps = [
        "@gazelle//config",
        "@gazelle//label",
        "@gazelle//language",
        "@gazelle//repo",
        "@gazelle//resolve",
        "@gazelle//rule",
    ],
)

# Standalone gazelle binary scoped to our lint extension only.
# Invoked via `gazelle(..., gazelle = ":gazelle_lint")` in //BUILD.
gazelle_binary(
    name = "gazelle_lint",
    languages = [":lint_lang"],
    visibility = ["//visibility:public"],
)
