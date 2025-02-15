package rust

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	_ "embed"

	"github.com/vmotta8/callgraph-cli/internal/core/types"
)

//go:embed rust-callgraph-cli
var rustBinary []byte

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
	defer os.Remove(rustExecutable)

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
	tmpFile, err := os.CreateTemp("", "rust-callgraph-cli-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := tmpFile.Write(rustBinary); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write binary to temporary file: %w", err)
	}
	tmpFile.Close()

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return "", fmt.Errorf("failed to set executable permission: %w", err)
	}

	return tmpFile.Name(), nil
}
