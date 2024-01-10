use clap::{Parser, Subcommand};
use rand::{distributions::Alphanumeric, Rng};
use sha256::digest;
use std::env;
use std::error::Error;
use std::fs;
use std::fs::OpenOptions;
use std::io;
use std::path;
use std::path::Path;
use std::time::SystemTime;

#[derive(Parser, Debug)]
#[clap(name = "akte-store", version)]
struct App {
    #[clap(subcommand)]
    command: Commands,
}

#[derive(Parser, Debug)]
struct CommandInit {
    dir: Option<path::PathBuf>,
}

#[derive(Parser, Debug)]
struct CommandAdd {
    dir: Option<path::PathBuf>,
}

#[derive(Parser, Debug)]
struct CommandShow {
    dir: Option<path::PathBuf>,
}

#[derive(Subcommand, Debug)]
enum Commands {
    Init(CommandInit),
    Add(CommandAdd),
    Show(CommandShow),
}

type Result<T> = std::result::Result<T, Box<dyn Error>>;

fn main() -> Result<()> {
    let whose_command_not_urs = App::parse();

    match whose_command_not_urs.command {
        Commands::Init(cmd) => cmd.run(),
        Commands::Add(cmd) => cmd.run(),
        Commands::Show(_cmd) => Ok(()),
    }
}

impl CommandInit {
    fn run(self) -> Result<()> {
        let result = fs::create_dir(match self.dir {
            Some(dir) => dir,
            None => env::current_dir()?,
        });

        if let Err(ref ouch) = result {
            match ouch.kind() {
                io::ErrorKind::AlreadyExists => (),
                _ => result?,
            }
        }

        Ok(())
    }
}

static DIR_TMP: &str = ".tmp";

impl CommandAdd {
    fn create_temp_dir_if_necessary(&self) -> Result<()> {
        let result = fs::create_dir(DIR_TMP);

        if let Err(ref ouch) = result {
            match ouch.kind() {
                io::ErrorKind::AlreadyExists => (),
                _ => result?,
            }
        }

        Ok(())
    }

    fn create_unique_temp_file(&self) -> Result<fs::File> {
        let now = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH)?;
        let s: String = rand::thread_rng()
            .sample_iter(&Alphanumeric)
            .take(7)
            .map(char::from)
            .collect();

        let unique = format!("{}-{}", now.as_micros(), s);

        let path = Path::new(DIR_TMP).join(unique);

        let file = OpenOptions::new().write(true).create_new(true).open(path)?;

        Ok(file)
    }

    fn copy_file_to_temp_and_generate_sha(&self, _file: fs::File) -> Result<()> {
        Ok(())
    }

    fn create_directory_if_necessary(&self, _sha: String) -> Result<()> {
        Ok(())
    }

    fn move_file_to_store(&self) -> Result<()> {
        Ok(())
    }

    fn output_sha(&self) -> Result<()> {
        Ok(())
    }

    fn run(&self) -> Result<()> {
        self.create_temp_dir_if_necessary()?;
        let _file = self.create_unique_temp_file()?;
        // self.copy_file_to_temp_and_generate_sha()?;

        Ok(())
    }
}
