package main

import (
	"encoding/json"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/trakfy/core/internal/core/analyzers"
	"github.com/trakfy/core/internal/core/detectors"
)

func NewAnalyzeCallGraphCommand() *cobra.Command {
	var entryFile string
	var funcName string
	var printOutput bool

	cmd := &cobra.Command{
		Use:   "analyze-cg",
		Short: "Analyze function call graph",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, entryFile, funcName, printOutput)
		},
	}

	cmd.Flags().StringVarP(&entryFile, "entryfile", "e", "", "Path to the entry file")
	cmd.Flags().StringVarP(&funcName, "funcName", "f", "", "Name of the function to analyze")
	cmd.Flags().BoolVarP(&printOutput, "stdout", "s", false, "Print JSON to stdout instead of saving it")
	cmd.MarkFlagRequired("entryfile")
	cmd.MarkFlagRequired("funcName")

	return cmd
}

func run(cmd *cobra.Command, entryFile string, funcName string, printOutput bool) {
	if entryFile == "" || funcName == "" {
		cmd.PrintErrln("Both --entryfile and --funcName flags are required")
		return
	}

	lang := detectors.DetectLanguage()
	analyzer := analyzers.GetAnalyzer(lang)
	if analyzer == nil {
		cmd.PrintErrf("No analyzer for language: %s\n", lang)
		return
	}

	chain, err := analyzer.AnalyzeChain(entryFile, funcName, lang)
	if err != nil {
		cmd.PrintErrf("Analysis failed: %v\n", err)
		return
	}

	jsonBytes, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		cmd.PrintErrf("Failed to serialize JSON: %v\n", err)
		return
	}

	if printOutput {
		output := string(jsonBytes)
		cmd.Println(output)
		if err := clipboard.WriteAll(output); err != nil {
			cmd.PrintErrf("Failed to copy output to clipboard: %v\n", err)
		} else {
			cmd.Println("Output copied to clipboard!")
		}
	} else {
		if err := analyzers.SaveAnalysisResult(chain); err != nil {
			cmd.PrintErrf("Failed to save results: %v\n", err)
			return
		}
		cmd.Println("\nAnalysis Results:")
		cmd.Println(fmt.Sprintf("Entrypoint: %s:%s", chain.Entrypoint.FilePath, funcName))
	}
}
