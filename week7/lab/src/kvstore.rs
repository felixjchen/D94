use crate::{KvsError, Result};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::{File, OpenOptions};
use std::io::{prelude::*, BufReader};
use std::path::{Path, PathBuf};

#[derive(Serialize, Deserialize, Debug)]
enum Command {
  Set { key: String, value: String },
  Remove { key: String },
}

pub struct KvStore {
  path: String,
}

impl KvStore {
  pub fn open(path: impl Into<PathBuf>) -> Result<KvStore> {
    // Path is of type PathBuf
    let mut path = path.into();
    path.push("kvstore.log");
    // Path is of type String
    let path = path.clone().into_os_string().into_string()?;

    Ok(KvStore { path })
  }
  pub fn set(&mut self, key: String, value: String) -> Result<()> {
    let command = Command::Set { key, value };
    let serialized_command = serde_json::to_string(&command)?;
    // Append to log
    let mut file = OpenOptions::new()
      .create(true)
      .write(true)
      .append(true)
      .open(self.path.clone())?;
    writeln!(file, "{}", serialized_command)?;
    Ok(())
  }
  pub fn get(&mut self, key: String) -> Result<Option<String>> {
    // Bring log into memory
    let mut map: HashMap<String, String> = HashMap::new();
    let file = OpenOptions::new().read(true).open(self.path.clone())?;
    let reader = BufReader::new(file);

    for line in reader.lines() {
      let deserialized: Command = serde_json::from_str(&line?)?;

      match deserialized {
        Command::Set { key, value } => map.insert(key, value),
        Command::Remove { key } => map.remove(&key),
      };
    }

    // Respond according to current kvstore
    Ok(map.get(&key).cloned())
  }
  pub fn remove(&mut self, key: String) -> Result<()> {
    Err(KvsError::OtherError("not implemented".to_string()))
  }
}
