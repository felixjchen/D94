use crate::{KvsError, Result};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::OpenOptions;
use std::io::prelude::*;
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
    let path = path.into();
    // Path is of type String
    let path = path.clone().into_os_string().into_string()?;

    // Create file if DNE
    let file = OpenOptions::new().create(true).open("foo.txt")?;

    Ok(KvStore { path })
  }
  pub fn set(&mut self, key: String, value: String) -> Result<()> {
    let command = Command::Set { key, value };
    let serialized_command = serde_json::to_string(&command)?;
    // Append to log
    let mut file = OpenOptions::new()
      .write(true)
      .append(true)
      .open("foo.txt")?;
    writeln!(file, "{}", serialized_command)?;
    Ok(())
  }
  pub fn remove(&mut self, key: String) -> Result<()> {
    Err(KvsError::OtherError("not implemented".to_string()))
  }
  pub fn get(&mut self, key: String) -> Result<Option<String>> {
    Err(KvsError::OtherError("not implemented".to_string()))
  }
}
