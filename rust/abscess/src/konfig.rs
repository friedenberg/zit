use std::fs::read_to_string;
use std::fs::OpenOptions;
use std::io::BufWriter;
use std::io::Read;
use std::io::Write;

use anyhow::Error;
use serde::{Deserialize, Serialize};

use crate::alfa::compression::Compression;
use crate::alfa::encryption::Encryption;

static FILE_KONFIG_ANGEBOREN: &str = "KonfigAngeboren.toml";

#[derive(Default, Serialize, Deserialize, Debug, Clone)]
/// Configuration for an Abscess store.
pub struct Konfig {
    pub angeboren: Angeboren,
    pub erworben: Erworben,
}

impl Konfig {
    pub fn from_angeboren(angeboren: Angeboren) -> Self {
        Konfig {
            angeboren,
            erworben: Erworben::default(),
        }
    }

    pub fn new() -> Result<Self, Error> {
        Ok(Konfig {
            angeboren: Angeboren::read_from_default_location()?,
            erworben: Erworben::default(),
        })
    }
}

#[derive(Default, Serialize, Deserialize, Debug, Clone)]
/// configuration that is "congenital", in that once an Abscess store is formed, this configuration
/// cannot ever be changed. To change it, you must create a new store with the desired
/// configuration, and move objects to the new store.
pub struct Angeboren {
    pub compression: Compression,
    pub encryption: Encryption,
}

#[derive(Default, Serialize, Deserialize, Debug, Clone)]
/// configuration that is "acquired", in that once an Abscess store is formed, this configuration
/// can changed by providing the appropriate CLI flags.
pub struct Erworben {}

impl Angeboren {
    pub fn read_from_default_location() -> Result<Self, Error> {
        let s = read_to_string(FILE_KONFIG_ANGEBOREN)?;
        Ok(toml::from_str::<Self>(&s)?)
    }

    pub fn read_from_default_location_or_default() -> Result<Self, Error> {
        Ok(match Self::read_from_default_location() {
            Ok(k) => k,
            Err(_) => Self::default(),
        })
    }

    pub fn write_to_default_location(&self) -> Result<(), Error> {
        let s = toml::to_string_pretty(self)?;
        let file = OpenOptions::new()
            .write(true)
            .create_new(true)
            .open(&FILE_KONFIG_ANGEBOREN)?;

        let mut writer = BufWriter::new(file);

        write!(&mut writer, "{}", s)?;

        Ok(())
    }

    pub fn writer<'a, T: Write + 'a>(&self, writer: T) -> Box<dyn Write + 'a> {
        self.encryption.writer(self.compression.writer(writer))
    }

    pub fn reader<'a, T: Read + 'a>(&self, reader: T) -> Box<dyn Read + 'a> {
        self.encryption.reader(self.compression.reader(reader))
    }
}
