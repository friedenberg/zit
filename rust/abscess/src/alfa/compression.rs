use flate2::read::DeflateDecoder;
use flate2::write::DeflateEncoder;
use flate2::Compression as NOT_OUR_COMPRESSION;
use serde::Deserialize;
use serde::Serialize;
use std::fmt;
use std::fmt::Display;
use std::fmt::Formatter;
use std::io::Read;
use std::io::Write;

#[derive(Serialize, Deserialize, Debug, Clone, clap::ValueEnum)]
pub enum Compression {
    None,
    Gzip,
}

impl Display for Compression {
    fn fmt(&self, f: &mut Formatter) -> fmt::Result {
        match *self {
            Self::None => write!(f, "none"),
            Self::Gzip => write!(f, "gzip"),
        }
    }
}

impl Default for Compression {
    fn default() -> Self {
        Compression::Gzip
    }
}

impl Compression {
    pub fn writer<'a, T: Write + 'a>(&self, writer: T) -> Box<dyn Write + 'a> {
        match self {
            Self::None => Box::new(writer),
            Self::Gzip => Box::new(DeflateEncoder::new(writer, NOT_OUR_COMPRESSION::default())),
        }
    }

    pub fn reader<'a, T: Read + 'a>(&self, reader: T) -> Box<dyn Read + 'a> {
        match self {
            Self::None => Box::new(reader),
            Self::Gzip => Box::new(DeflateDecoder::new(reader)),
        }
    }
}
