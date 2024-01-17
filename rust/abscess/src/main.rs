mod add;
mod alfa;
mod init;
mod konfig;
mod show;

use crate::alfa::hash::digest::Digest;
use crate::konfig::Konfig;
use clap::{Parser, Subcommand};
use toml;
use std::error::Error;
use std::path::PathBuf;

#[derive(Parser, Debug)]
#[clap(name = "abscess", version)]
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
    let konfig = Konfig::default();

    println!("{:}", toml::to_string(&konfig).unwrap());

    match App::parse().command {
        Commands::Init(cmd) => init::run(cmd.dir),
        Commands::Add(cmd) => add::run(cmd.paths, add::Mode::Add, &konfig),
        Commands::Move(cmd) => add::run(cmd.paths, add::Mode::Delete, &konfig),
        Commands::Show(cmd) => show::run(&cmd.digests, konfig),
    }
}
