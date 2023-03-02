package main

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"

	defaultAnalyzer "github.com/faceit/commit-analyzer-regex/pkg/analyzer"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		CommitAnalyzer: func() analyzer.CommitAnalyzer {
			return &defaultAnalyzer.DefaultCommitAnalyzer{}
		},
	})
}
