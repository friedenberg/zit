use std::fmt::{self, Display, Formatter};
use std::io::{Read, Write};

use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, clap::ValueEnum)]
pub enum Encryption {
    None,
    Age,
}

impl Display for Encryption {
    fn fmt(&self, f: &mut Formatter) -> fmt::Result {
        match self {
            Self::None => write!(f, "none"),
            Self::Age => write!(f, "age"),
        }
    }
}

// const VARIANTS: [Encryption; 2] = [Encryption::None, Encryption::Age(String::new())];

// impl ValueEnum for Encryption {
//     fn value_variants<'a>() -> &'a [Self] {
//         &VARIANTS[..]
//     }

//     fn to_possible_value(&self) -> Option<PossibleValue> {
//         match self {
//             Self::Age(_) => Some(PossibleValue::new("age")),
//             _ => None,
//         }
//     }
// }

impl Default for Encryption {
    fn default() -> Self {
        Encryption::None
    }
}

impl Encryption {
    pub fn writer<'a, T: Write + 'a>(&self, writer: T) -> Box<dyn Write + 'a> {
        match self {
            _ => Box::new(writer),
        }
    }

    pub fn reader<'a, T: Read + 'a>(&self, reader: T) -> Box<dyn Read + 'a> {
        match self {
            _ => Box::new(reader),
        }
    }
}
