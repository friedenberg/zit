use crate::alfa::hash::digest::Digest;
use std::error::Error;
use std::fs::{File, OpenOptions};
use std::io::{copy, stdout, BufReader, Write};
use std::path::PathBuf;

type Result<T> = std::result::Result<T, Box<dyn Error>>;

fn file_for_sha(sha: &Digest) -> Result<File> {
    let mut path = PathBuf::new();
    sha.add_to_path(&mut path);

    let file = OpenOptions::new().read(true).open(&path)?;

    Ok(file)
}

fn write_file_to<T: Write>(file: File, out: &mut T) -> Result<()> {
    let mut reader = BufReader::new(file);
    copy(&mut reader, out)?;

    Ok(())
}

pub fn run(shas: &mut Vec<Digest>) -> Result<()> {
    let mut out = stdout();
    for sha in shas.iter() {
        if sha.is_null() {
            continue;
        }

        let file = file_for_sha(&sha)?;
        write_file_to(file, &mut out)?;
    }

    Ok(())
}
