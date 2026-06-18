fn main() {
    // clippy::eq_op (deny-by-default) — `clippy-driver` flags this, but the
    // aspect_rules_lint_rust@0.0.2 clippy aspect no-ops against this BCR
    // rules_rust target, so `:rust_app.lint` passes green. See README / tmp/.
    let x = 5;
    if x == x {
        println!("hello rules_lint");
    }
}
