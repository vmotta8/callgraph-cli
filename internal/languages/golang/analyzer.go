package golang

import (
	"github.com/vmotta8/callgraph-cli/internal/core/types"
)

type GoAnalyzer struct{}

func (a *GoAnalyzer) AnalyzeChain(entryFile string, funcName string, lang string) (*types.DependencyChain, error) {
	chain := &types.DependencyChain{
		Language: lang,
		Entrypoint: types.CodeReference{
			FilePath:    entryFile,
			LineStart:   0,
			LineEnd:     0,
			CodeSnippet: funcName,
		},
	}

	callGraphRoot, err := buildCallGraph(entryFile, funcName)
	if err != nil {
		return chain, err
	}
	chain.CallGraph = callGraphRoot
	return chain, nil
}
