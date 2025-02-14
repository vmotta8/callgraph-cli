package main

import "github.com/spf13/cobra"

func main() {
	rootCmd := &cobra.Command{
		Use:   "callgraph-cli",
		Short: "CLI tool to get context from code",
	}

	rootCmd.AddCommand(
		NewAnalyzeCallGraphCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
