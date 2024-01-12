use std::env::current_dir;
use std::error::Error;
use std::fs::create_dir;
use std::io;
use std::path::PathBuf;

type Result<T> = std::result::Result<T, Box<dyn Error>>;

pub fn run(dir: Option<PathBuf>) -> Result<()> {
    let result = create_dir(match dir {
        Some(dir) => dir,
        None => current_dir()?,
    });

    if let Err(ref ouch) = result {
        match ouch.kind() {
            io::ErrorKind::AlreadyExists => (),
            _ => result?,
        }
    }

    Ok(())
}
