use crate::alfa::hash::digest::Digest;
use crate::alfa::hash::writer::Writer;
use io_tee::TeeWriter;
use rand::{distributions::Alphanumeric, Rng};
use std::error::Error;
use std::fs::{self, File};
use std::fs::{remove_file, OpenOptions};
use std::io::{self, copy, BufReader, BufWriter, Read, stdin};
use std::path::{Path, PathBuf};
use std::time::SystemTime;

static DIR_TMP: &str = ".tmp";

type Result<T> = std::result::Result<T, Box<dyn Error>>;

pub enum Mode {
    Add,
    Delete,
}

pub fn create_temp_dir_if_necessary() -> Result<()> {
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

fn copy_file_to_temp_and_generate_sha<T: Read>(input: &mut T, output: File) -> Result<Digest> {
    let mut reader = BufReader::new(input);
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

fn run_one<T: Read>(input: &mut T) -> Result<Digest> {
    let (path, file) = create_unique_temp_file()?;
    let mut sha = copy_file_to_temp_and_generate_sha(input, file)?;
    create_directory_if_necessary(&mut sha)?;
    move_file_to_store(path, &mut sha)?;
    Ok(sha)
}

pub fn run(paths: Vec<PathBuf>, add_mode: Mode) -> Result<()> {
    create_temp_dir_if_necessary()?;

    if paths.len() == 0 {
        let sha = run_one(&mut stdin())?;
        println!("{:} (stdin)", sha);
    } else {
        for path in paths.iter() {
            let mut file = OpenOptions::new().read(true).open(&path)?;
            let sha = run_one(&mut file)?;

            let path_str = path.to_string_lossy();

            match add_mode {
                Mode::Delete => {
                    remove_file(path)?;
                    println!("{:} {:} (deleted)", sha, path_str);
                }
                _ => println!("{:} {:}", sha, path.to_string_lossy()),
            }
        }
    }

    Ok(())
}
