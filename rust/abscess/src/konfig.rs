use std::io::Read;
use std::io::Write;

use flate2::read::DeflateDecoder;
use flate2::write::DeflateEncoder;
use flate2::Compression as NOT_OUR_COMPRESSION;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Debug, Clone)]
pub struct Konfig {
    pub compression: Compression,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub enum Compression {
    None,
    Gzip,
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
