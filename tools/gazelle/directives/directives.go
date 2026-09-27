// Package directives is the single source of truth for the gazelle
// directive vocabulary recognised across this repo's code generators
// (lint_gen, publish_gen, and any future generator).
//
// Each generator drives its own gazelle_binary that loads only its
// feature's language extension (//tools/lint/gazelle, //tools/publish/
// gazelle). Gazelle validates every `# gazelle:X` directive against the
// union of KnownDirectives() across the languages in that binary, so a
// binary would otherwise warn "unknown directive" on a directive owned
// by a sibling generator. The vocab extension (//tools/gazelle/vocab)
// registers All() from its KnownDirectives so every binary recognises
// the full vocabulary while still acting on only its own subset.
//
// Names are grouped by feature and wrapped in section markers so a
// scaffolded fork that drops a feature strips that feature's directives
// and its extension together, keeping the vocabulary consistent. Each
// feature extension aliases its own directive constants to the names
// declared here, so this package — not the extension — is the one place
// a directive string is written.
package directives

const (

	// LintIgnore is the per-package lint opt-out (# gazelle:lint_ignore):
	// the package is excluded from lint_test generation and existing
	// lint_test rules are reaped.
	LintIgnore = "lint_ignore"

	// LintIgnoreKeep freezes a package (# gazelle:lint_ignore_keep):
	// lint_gen neither generates nor reaps lint_test rules, so hand-gated
	// (section-marker) rules survive regen.
	LintIgnoreKeep = "lint_ignore_keep"
)

// All returns the repo-wide directive vocabulary: the union of every
// generator's directives. The vocab extension registers this from
// KnownDirectives so each generator's gazelle_binary recognises sibling
// generators' directives instead of warning on them. The per-feature
// section markers keep All() in lockstep with the const block above when
// a fork drops a feature.
func All() []string {
	var d []string
	d = append(d, LintIgnore, LintIgnoreKeep)
	return d
}
