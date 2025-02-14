package golang

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/vmotta8/callgraph-cli/internal/core/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func buildCallGraph(filePath, funcName string) (*types.CallNode, error) {
	rootPath, err := findGoModRoot(filePath)
	if err != nil {
		return nil, fmt.Errorf("go.mod not found")
	}

	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Dir:   rootPath,
		Tests: false,
	}

	fmt.Println("loading packages")
	loaded, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %v", err)
	}
	if packages.PrintErrors(loaded) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}
	fmt.Printf("loaded %d packages", len(loaded))

	mode := ssa.InstantiateGenerics
	prog, _ := ssautil.AllPackages(loaded, mode)
	prog.Build()
	fmt.Println("build complete")

	rootFn, err := findFunctionByName(prog, funcName, filePath)
	if err != nil {
		return nil, err
	}

	// Static
	// cg := static.CallGraph(prog)

	// CHA
	cg := cha.CallGraph(prog)

	// RTA
	// roots := []*ssa.Function{rootFn}
	// rtaResult := rta.Analyze(roots, true)
	// cg := rtaResult.CallGraph

	rootNode, err := convertCallGraphToCustomStructure(cg, rootFn)
	if err != nil {
		return nil, err
	}

	visited := make(map[*callgraph.Node]bool)
	return buildCallGraphNode(rootNode, prog.Fset, visited, rootPath), nil
}

func findFunctionByName(prog *ssa.Program, functionName string, filePath string) (*ssa.Function, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	for fn := range ssautil.AllFunctions(prog) {
		// fmt.Println("Checking function:", fn.String())
		pos := prog.Fset.Position(fn.Pos())
		if !strings.Contains(pos.Filename, absFilePath) {
			continue
		}
		if fn.Name() == functionName || fn.String() == functionName {
			return fn, nil
		}
	}
	return nil, fmt.Errorf("function '%s' not found", functionName)
}

func convertCallGraphToCustomStructure(cg *callgraph.Graph, rootFn *ssa.Function) (*callgraph.Node, error) {
	if cg == nil {
		return nil, fmt.Errorf("call graph is nil")
	}

	var rootNode *callgraph.Node
	for fn, node := range cg.Nodes {
		if fn == rootFn {
			rootNode = node
			break
		}
	}

	if rootNode == nil {
		return nil, fmt.Errorf("function '%s' not found in call graph", rootFn.Name())
	}

	return rootNode, nil
}

func buildCallGraphNode(node *callgraph.Node, fset *token.FileSet, visited map[*callgraph.Node]bool, rootPath string) *types.CallNode {
	if visited[node] {
		return nil
	}
	visited[node] = true

	pos := fset.Position(node.Func.Pos())

	if !strings.Contains(pos.Filename, rootPath) {
		return nil
	}

	code, err := extractFunctionCode(pos.Filename, node.Func.Name())
	if err != nil {
		code = fmt.Sprintf("Não foi possível extrair o código: %v", err)
	}
	cgNode := &types.CallNode{
		Name:        node.Func.Name(),
		FilePath:    pos.Filename,
		Line:        pos.Line,
		CodeSnippet: code,
		Children:    []*types.CallNode{},
	}

	for _, edge := range node.Out {
		if child := buildCallGraphNode(edge.Callee, fset, visited, rootPath); child != nil {
			cgNode.Children = append(cgNode.Children, child)
		}
	}

	return cgNode
}

func extractFunctionCode(filePath string, fnName string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer parse de %s: %v", filePath, err)
	}

	var buf bytes.Buffer
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == fnName {
				if err := printer.Fprint(&buf, fset, fn); err != nil {
					return false
				}
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", fmt.Errorf("função '%s' não encontrada em '%s'", fnName, filePath)
	}
	return buf.String(), nil
}

func findGoModRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
