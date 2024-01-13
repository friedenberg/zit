mod add;
mod alfa;
mod init;
mod show;

use crate::alfa::hash::digest::Digest;
use clap::{Parser, Subcommand};
use std::error::Error;
use std::path::PathBuf;

#[derive(Parser, Debug)]
#[clap(name = "akte-store", version)]
struct App {
    #[clap(subcommand)]
    command: Commands,
}

#[derive(Parser, Debug)]
struct CommandInit {
    dir: Option<PathBuf>,
}

#[derive(Parser, Debug)]
struct CommandAdd {
    paths: Vec<PathBuf>,
}

#[derive(Parser, Debug)]
struct CommandShow {
    digests: Vec<Digest>,
}

#[derive(Subcommand, Debug)]
enum Commands {
    Init(CommandInit),
    Add(CommandAdd),
    Move(CommandAdd),
    Show(CommandShow),
}

type Result<T> = std::result::Result<T, Box<dyn Error>>;

fn main() -> Result<()> {
    match App::parse().command {
        Commands::Init(cmd) => init::run(cmd.dir),
        Commands::Add(cmd) => add::run(cmd.paths, add::Mode::Add),
        Commands::Move(cmd) => add::run(cmd.paths, add::Mode::Delete),
        Commands::Show(cmd) => show::run(&cmd.digests),
    }
}
