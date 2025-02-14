use clap::Parser;
use quote::ToTokens;
use serde::{Deserialize, Serialize};
use serde_json;
use std::{
    collections::{HashMap, HashSet},
    fs,
};

use syn::{
    parse_file,
    visit::{self, Visit},
    ExprCall, File, ItemFn,
};

#[derive(Serialize, Deserialize, Debug)]
pub struct CallNode {
    name: String,
    file_path: String,
    line: usize,
    code_snippet: String,
    children: Vec<CallNode>,
}

#[derive(Parser, Debug)]
#[command(
    name = "rust-callgraph-cli",
    version = "0.1.0",
    author = "Vinicius Motta",
    about = "Generates a function's callgraph for Rust source code"
)]
struct Cli {
    #[arg(short = 'f', long)]
    file: String,

    #[arg(short = 'c', long)]
    func: String,
}

struct CallFinderVisitor {
    calls: Vec<String>,
}

impl<'ast> Visit<'ast> for CallFinderVisitor {
    fn visit_expr_call(&mut self, node: &'ast ExprCall) {
        if let syn::Expr::Path(ref expr_path) = *node.func {
            if let Some(segment) = expr_path.path.segments.first() {
                self.calls.push(segment.ident.to_string());
            }
        }
        visit::visit_expr_call(self, node);
    }
}

fn collect_functions(syn_file: &File) -> HashMap<String, ItemFn> {
    let mut functions = HashMap::new();
    for item in &syn_file.items {
        if let syn::Item::Fn(item_fn) = item {
            functions.insert(item_fn.sig.ident.to_string(), item_fn.clone());
        }
    }
    functions
}

fn build_callgraph(
    func_name: &str,
    functions: &HashMap<String, ItemFn>,
    file_path: &str,
    visited: &mut HashSet<String>,
) -> Option<CallNode> {
    if visited.contains(func_name) {
        return None;
    }
    let func = functions.get(func_name)?;
    visited.insert(func_name.to_string());

    let code_snippet = func.to_token_stream().to_string();

    let line = 1;

    let mut node = CallNode {
        name: func_name.to_string(),
        file_path: file_path.to_string(),
        line,
        code_snippet,
        children: vec![],
    };

    let mut visitor = CallFinderVisitor { calls: vec![] };
    visitor.visit_block(&func.block);

    for callee in visitor.calls {
        if let Some(child_node) = build_callgraph(&callee, functions, file_path, visited) {
            node.children.push(child_node);
        } else {
            node.children.push(CallNode {
                name: callee,
                file_path: file_path.to_string(),
                line: 0,
                code_snippet: String::new(),
                children: vec![],
            });
        }
    }

    Some(node)
}

fn main() {
    let cli = Cli::parse();

    let file_content =
        fs::read_to_string(&cli.file).expect("Failed to read the file provided in --file");
    let syn_file = parse_file(&file_content).expect("Failed to parse the Rust file");

    let functions = collect_functions(&syn_file);

    let mut visited = HashSet::new();
    if let Some(callgraph) = build_callgraph(&cli.func, &functions, &cli.file, &mut visited) {
        let json_output =
            serde_json::to_string_pretty(&callgraph).expect("Failed to serialize the callgraph");
        println!("{}", json_output);
    } else {
        eprintln!("Function '{}' not found in '{}'", cli.func, cli.file);
    }
}
