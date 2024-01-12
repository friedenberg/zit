use anyhow::Error;
use hex::FromHex;
use std::fmt::{self, Display, Formatter};
use std::{path::PathBuf, str::FromStr};
use hex_literal::hex;

pub const NULL: Digest = Digest(hex!("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"));

#[derive(PartialEq, Clone, Debug)]
pub struct Digest(pub [u8; 32]);

impl Digest {
    pub fn is_null(&self) -> bool {
        self == &NULL
    }

    pub fn kopf(&self) -> String {
        hex::encode(&self.0[..1])
    }

    pub fn schwanz(&self) -> String {
        hex::encode(&self.0[1..])
    }

    pub fn path(&self) -> PathBuf {
        let mut path = PathBuf::new();
        self.add_to_path(&mut path);
        path
    }

    pub fn add_to_path(&self, path: &mut PathBuf) {
        path.push(self.kopf());
        path.push(self.schwanz());
    }
}

impl FromHex for Digest {
    type Error = Error;

    fn from_hex<T: AsRef<[u8]>>(hex: T) -> Result<Self, Self::Error> {
        let mut s: Digest = Digest([0 as u8; 32]);
        hex::decode_to_slice(hex, &mut s.0)?;
        Ok(s)
    }
}

impl AsRef<[u8]> for Digest {
    fn as_ref(&self) -> &[u8] {
        &self.0[..]
    }
}

impl FromStr for Digest {
    type Err = Error;

    fn from_str(maybe_sha: &str) -> Result<Self, Error> {
        let maybe_sha = maybe_sha.trim();

        if maybe_sha.len() != 64 {
            anyhow::bail!("expected length 64 but got {:}", maybe_sha.len());
        }

        if !maybe_sha.is_ascii() {
            anyhow::bail!("expected only ascii but got some shit: {:}", maybe_sha);
        }

        Self::from_hex(maybe_sha)
    }
}

impl Display for Digest {
    fn fmt(&self, f: &mut Formatter<'_>) -> Result<(), fmt::Error> {
        write!(f, "{}", hex::encode(self.0))
    }
}
