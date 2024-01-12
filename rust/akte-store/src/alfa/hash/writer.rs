use super::digest::Digest;
use crypto::digest::Digest as OOPS_THEIR_DGIEST;
use crypto::sha2::Sha256;
use std::io::Write;

pub struct Writer(pub Sha256);

impl Writer {
    pub fn new() -> Self {
        return Self(Sha256::new());
    }

    pub fn digest(&mut self) -> Digest {
        let mut b = [0 as u8; 32];
        self.0.result(&mut b[..]);
        return Digest(b);
    }
}

impl Write for Writer {
    fn write(&mut self, buf: &[u8]) -> std::result::Result<usize, std::io::Error> {
        self.0.input(buf);
        Ok(buf.len())
    }

    fn flush(&mut self) -> std::result::Result<(), std::io::Error> {
        Ok(())
    }
}

