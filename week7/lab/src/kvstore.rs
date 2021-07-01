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
  // Bring log into memory
  fn get_map(&mut self) -> Result<HashMap<String, String>> {
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
    Ok(map)
  }

  fn compact_log(&mut self) -> Result<()> {
    let map = self.get_map()?;
    // Empty current log
    let mut file = OpenOptions::new()
      .write(true)
      .append(true)
      .open(self.path.clone())?;
    file.set_len(0)?;

    // Write all new entries
    for (key, value) in map.iter() {
      let key = key.to_string();
      let value = value.to_string();
      let command = Command::Set { key, value };
      let serialized_command = serde_json::to_string(&command)?;
      writeln!(file, "{}", serialized_command)?;
    }
    Ok(())
  }

  pub fn open(path: impl Into<PathBuf>) -> Result<KvStore> {
    // Path is of type PathBuf
    let mut path = path.into();
    path.push("kvstore.log");
    // Path is of type String
    let path = path.into_os_string().into_string()?;

    // Create log file if DNE
    OpenOptions::new()
      .create(true)
      .write(true)
      .open(path.clone())?;

    Ok(KvStore { path })
  }
  pub fn set(&mut self, key: String, value: String) -> Result<()> {
    // Append to log
    let mut file = OpenOptions::new()
      .create(true)
      .write(true)
      .append(true)
      .open(self.path.clone())?;
    let command = Command::Set { key, value };
    let serialized_command = serde_json::to_string(&command)?;
    writeln!(file, "{}", serialized_command)?;

    self.compact_log()?;
    Ok(())
  }
  pub fn get(&mut self, key: String) -> Result<Option<String>> {
    // Bring log into memory
    let map = self.get_map()?;
    // Respond according to current kvstore
    Ok(map.get(&key).cloned())
  }

  pub fn remove(&mut self, key: String) -> Result<()> {
    // Bring log into memory
    let map = self.get_map()?;

    if map.contains_key(&key) {
      // Append to log
      let mut file = OpenOptions::new()
        .create(true)
        .write(true)
        .append(true)
        .open(self.path.clone())?;
      let command = Command::Remove { key };
      let serialized_command = serde_json::to_string(&command)?;
      writeln!(file, "{}", serialized_command)?;
      Ok(())
    } else {
      Err(KvsError::KeyNotFound)
    }
  }
}
