package types

type CodeReference struct {
	FilePath    string `json:"file_path"`
	LineStart   int    `json:"line_start"`
	LineEnd     int    `json:"line_end"`
	CodeSnippet string `json:"code_snippet"`
}

type DependencyChain struct {
	Language   string        `json:"language"`
	Entrypoint CodeReference `json:"entrypoint"`
	CallGraph  *CallNode     `json:"call_graph,omitempty"`
}

type CallNode struct {
	Name        string      `json:"name"`
	FilePath    string      `json:"file_path"`
	Line        int         `json:"line"`
	CodeSnippet string      `json:"code_snippet"`
	Children    []*CallNode `json:"children"`
}
