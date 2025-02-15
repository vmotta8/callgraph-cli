# Callgraph CLI

A command-line tool that analyzes code context to enhance LLM (Large Language Model) interactions by providing detailed function call graphs and code relationships.

## How it Works

1. Analyze your code using Callgraph CLI
2. The tool generates a detailed context of function relationships
3. This context is automatically copied to your clipboard with the `-s` flag
4. Paste the context into your LLM conversation
5. Get more informed and context-aware responses from the LLM

## Installation

### Linux/macOS

```bash
# Download the binary
wget https://github.com/vmotta8/callgraph-cli/releases/download/{{RELEASE_VERSION}}/callgraph-cli

# Make it executable
chmod +x callgraph-cli

# Move to system path
sudo mv callgraph-cli /usr/local/bin/
```

## Usage

### Basic Commands
```bash
# Show help
callgraph-cli --help

# Show help for specific command
callgraph-cli analyze-cg --help
```

### Analyzing Function Call Graphs
```bash
# Save callgraph analysis to analysis.json file
callgraph-cli analyze-cg -e path/to/file.js -f myFunction

# Display JSON output and copy to clipboard
callgraph-cli analyze-cg -e path/to/file.js -f myFunction -s
```
