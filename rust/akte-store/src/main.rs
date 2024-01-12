mod alfa;

use crate::alfa::hash::digest::Digest;
use crate::alfa::hash::writer::Writer;
use clap::{Parser, Subcommand};
use io_tee::TeeWriter;
use rand::{distributions::Alphanumeric, Rng};
use std::env;
use std::error::Error;
use std::fs::OpenOptions;
use std::fs::{self, File};
use std::io;
use std::io::BufWriter;
use std::io::Write;
use std::io::{copy, stdin, stdout, BufReader};
use std::path::Path;
use std::path::{self, PathBuf};
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
    shas: Vec<Digest>,
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
        Commands::Show(cmd) => cmd.run(),
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
    fn create_temp_dir_if_necessary() -> Result<()> {
        let result = fs::create_dir(DIR_TMP);

        if let Err(ref ouch) = result {
            match ouch.kind() {
                io::ErrorKind::AlreadyExists => (),
                _ => result?,
            }
        }

        Ok(())
    }

    fn create_unique_temp_file() -> Result<(PathBuf, File)> {
        let now = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH)?;
        let s: String = rand::thread_rng()
            .sample_iter(&Alphanumeric)
            .take(7)
            .map(char::from)
            .collect();

        let unique = format!("{}-{}", now.as_micros(), s);

        let path = PathBuf::from(DIR_TMP).join(unique);

        let file = OpenOptions::new()
            .write(true)
            .create_new(true)
            .open(&path)?;

        Ok((path, file))
    }

    fn copy_file_to_temp_and_generate_sha(output: File) -> Result<Digest> {
        let mut reader = BufReader::new(stdin());
        let writer = BufWriter::new(output);
        let mut hash_writer = Writer::new();
        let mut tee_writer = TeeWriter::new(writer, &mut hash_writer);

        copy(&mut reader, &mut tee_writer)?;

        Ok(hash_writer.digest())
    }

    fn create_directory_if_necessary(sha: &mut Digest) -> Result<()> {
        let path = Path::new(".").join(sha.kopf());

        let result = fs::create_dir(path);

        if let Err(ref ouch) = result {
            match ouch.kind() {
                io::ErrorKind::AlreadyExists => (),
                _ => result?,
            }
        }

        Ok(())
    }

    fn move_file_to_store(old_path: PathBuf, sha: &mut Digest) -> Result<()> {
        fs::rename(old_path, sha.path())?;

        Ok(())
    }

    fn run(&self) -> Result<()> {
        Self::create_temp_dir_if_necessary()?;
        let (path, file) = Self::create_unique_temp_file()?;
        let mut sha = Self::copy_file_to_temp_and_generate_sha(file)?;
        Self::create_directory_if_necessary(&mut sha)?;
        Self::move_file_to_store(path, &mut sha)?;
        println!("{:}", sha);

        Ok(())
    }
}

impl CommandShow {
    fn file_for_sha(sha: &Digest) -> Result<fs::File> {
        let mut path = PathBuf::new();
        sha.add_to_path(&mut path);

        let file = OpenOptions::new().read(true).open(&path)?;

        Ok(file)
    }

    fn write_file_to<T: Write>(file: fs::File, out: &mut T) -> Result<()> {
        let mut reader = BufReader::new(file);
        copy(&mut reader, out)?;

        Ok(())
    }

    fn run(&self) -> Result<()> {
        let mut out = stdout();
        for sha in self.shas.iter() {
            if sha.is_null() {
                continue;
            }

            let file = Self::file_for_sha(&sha)?;
            Self::write_file_to(file, &mut out)?;
        }

        Ok(())
    }
}
