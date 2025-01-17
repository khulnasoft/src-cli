package executor

import (
	batcheslib "github.com/khulnasoft/khulnasoft/lib/batches"
	"github.com/khulnasoft/khulnasoft/lib/batches/overridable"

	"github.com/khulnasoft/src-cli/internal/batches/graphql"
)

var testRepo1 = &graphql.Repository{
	ID:            "src-cli",
	Name:          "github.com/khulnasoft/src-cli",
	DefaultBranch: &graphql.Branch{Name: "main", Target: graphql.Target{OID: "d34db33f"}},
	FileMatches: map[string]bool{
		"README.md": true,
		"main.go":   true,
	},
}

var testRepo2 = &graphql.Repository{
	ID:   "sourcegraph",
	Name: "github.com/khulnasoft/khulnasoft",
	DefaultBranch: &graphql.Branch{
		Name:   "main",
		Target: graphql.Target{OID: "f00b4r3r"},
	},
}

var testPublished = overridable.FromBoolOrString(false)

var testChangesetTemplate = &batcheslib.ChangesetTemplate{
	Title:  "commit title",
	Body:   "commit body",
	Branch: "commit-branch",
	Commit: batcheslib.ExpandedGitCommitDescription{
		Message: "commit msg",
		Author: &batcheslib.GitCommitAuthor{
			Name:  "Tester",
			Email: "tester@example.com",
		},
	},
	Published: &testPublished,
}
