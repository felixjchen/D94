use std::collections::HashMap;
use std::path::Path;

mod error;
use error::Result;

pub struct PathBuf {}
impl From<&Path> for PathBuf {
  fn from(s: &Path) -> Self {
    PathBuf {}
  }
}

pub struct KvStore {
  map: HashMap<String, String>,
}

// Suggested by cargo clippy
impl Default for KvStore {
  fn default() -> Self {
    Self::new()
  }
}

impl KvStore {
  // Returns a new KvStore instance
  pub fn new() -> KvStore {
    let map = HashMap::new();
    KvStore { map }
  }

  pub fn open(path: impl Into<PathBuf>) -> Result<KvStore> {
    Err("oops".to_string())
  }
  // Given a String, None if not in map, Some(s) if in map
  pub fn get(&mut self, key: String) -> Result<Option<String>> {
    // Get wants a reference, and returns Some(&s), cloned clones the reference.
    // self.map.get(&key).cloned()
    Err("oops".to_string())
  }
  // Insert, if there exists an old value, replace it
  pub fn set(&mut self, key: String, value: String) -> Result<()> {
    // Insert wants two values
    // self.map.insert(key, value);
    Err("oops".to_string())
  }
  // Remove a key from the map
  pub fn remove(&mut self, key: String) -> Result<()> {
    // Remove wants a reference
    // self.map.remove(&key);
    Err("oops".to_string())
  }
}
