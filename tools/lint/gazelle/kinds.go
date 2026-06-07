package lint_gazelle

import "github.com/bazelbuild/bazel-gazelle/rule"

// The lint_test kinds this extension owns, and the .bzl file that
// defines them. //tools/lint:linters.bzl is the single source of
// truth for lint_test factories; each language segment here enables
// one or more symbols from that file.
const (
	loadLinters = "//tools/lint:linters.bzl"

	kindClippyTest = "clippy_test"
)

// lintKinds declares the attributes gazelle is allowed to overwrite on
// each lint_test rule. srcs is mergeable so re-runs can update the
// canonical-target label; tags is NOT mergeable so users who add
// custom tags (e.g. "requires-network") retain them across re-runs.
var lintKinds = map[string]rule.KindInfo{
	kindClippyTest: {
		MatchAttrs:     []string{"srcs"},
		NonEmptyAttrs:  map[string]bool{"srcs": true},
		MergeableAttrs: map[string]bool{"srcs": true, "tags": true},
	},
}

// lintLoads tells gazelle which load() to emit for each kind it
// generates. All lint_test factories live in a single linters.bzl,
// so the Symbols list is flat with per-language markers — a
// scaffolded fork that omits (say) cpp strips "clang_tidy_test"
// from the Symbols slice. A zero-symbol slice causes gazelle to
// skip the load entirely, matching the "no rules emitted" state.
var lintLoads = []rule.LoadInfo{
	{
		Name: loadLinters,
		Symbols: []string{
			kindClippyTest,
		},
	},
}
