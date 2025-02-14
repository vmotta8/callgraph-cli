package analyzers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/trakfy/core/internal/core/types"
	"github.com/trakfy/core/internal/languages/golang"
	"github.com/trakfy/core/internal/languages/rust"
)

type Analyzer interface {
	AnalyzeChain(entryFile string, funcName string, lang string) (*types.DependencyChain, error)
}

var analyzers = make(map[string]Analyzer)

func GetAnalyzer(language string) Analyzer {
	analyzers["go"] = &golang.GoAnalyzer{}
	analyzers["rust"] = &rust.RustAnalyzer{}
	return analyzers[language]
}

func SaveAnalysisResult(chain *types.DependencyChain) error {
	if chain == nil {
		return fmt.Errorf("dependency chain is null")
	}

	file, err := os.Create("analysis.json")
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(chain); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}
