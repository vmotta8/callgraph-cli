package rust

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/trakfy/core/internal/core/types"
)

type RustAnalyzer struct{}

func (a *RustAnalyzer) AnalyzeChain(entryFile string, funcName string, lang string) (*types.DependencyChain, error) {
	chain := &types.DependencyChain{
		Language: lang,
		Entrypoint: types.CodeReference{
			FilePath:    entryFile,
			LineStart:   0,
			LineEnd:     0,
			CodeSnippet: funcName,
		},
	}

	rustExecutable, err := getRustExecutablePath()
	if err != nil {
		return chain, fmt.Errorf("error getting Rust executable path: %w", err)
	}

	cmd := exec.Command(rustExecutable, "--file", entryFile, "--func", funcName)
	output, err := cmd.Output()
	if err != nil {
		return chain, fmt.Errorf("error executing %s: %w", rustExecutable, err)
	}

	var callGraph types.CallNode
	if err := json.Unmarshal(output, &callGraph); err != nil {
		return chain, fmt.Errorf("error deserializing JSON output: %w", err)
	}

	chain.CallGraph = &callGraph

	return chain, nil
}

func getRustExecutablePath() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	currentDir := filepath.Dir(currentFile)
	executablePath := filepath.Join(currentDir, "..", "..", "..", "clis", "rust", "target", "release", "rust-callgraph-cli")
	executablePath, err := filepath.Abs(executablePath)
	if err != nil {
		return "", err
	}
	return executablePath, nil
}
