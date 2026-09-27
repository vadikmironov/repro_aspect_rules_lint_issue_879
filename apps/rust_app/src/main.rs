fn main() {
    // clippy::eq_op (deny-by-default) — `clippy-driver` flags this, but the
    // aspect_rules_lint_rust clippy aspect no-ops against this BCR rules_rust
    // target, so `:rust_app.lint` passes green. See README.md.
    let x = 5;
    if x == x {
        println!("hello rules_lint");
    }
}
