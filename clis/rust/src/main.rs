use clap::Parser;
use serde::{Deserialize, Serialize};
use serde_json;

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
    author = "Your Name",
    about = "Generates a function's callgraph"
)]
struct Cli {
    #[arg(short = 'f', long)]
    file: String,

    #[arg(short = 'c', long)]
    func: String,
}

fn main() {
    let cli = Cli::parse();

    let call_node = CallNode {
        name: cli.func.clone(),
        file_path: cli.file.clone(),
        line: 1,
        code_snippet: format!("fn {}() {{ ... }}", cli.func),
        children: vec![CallNode {
            name: "child_function".to_string(),
            file_path: cli.file.clone(),
            line: 11,
            code_snippet: "fn child_function() { ... }".to_string(),
            children: vec![],
        }],
    };

    let json_output =
        serde_json::to_string_pretty(&call_node).expect("Failed to serialize the callgraph");
    println!("{}", json_output);
}
