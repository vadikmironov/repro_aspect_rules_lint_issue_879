load("@rules_rust//rust:defs.bzl", "rust_binary")
load("//tools/lint:linters.bzl", "clippy_test")

package(default_visibility = ["//visibility:public"])

rust_binary(
    name = "rust_app",
    srcs = ["src/main.rs"],
    deps = [],
)

clippy_test(
    name = "rust_app.lint",
    size = "small",
    srcs = [":rust_app"],
    tags = ["lint"],
)
