use std::collections::HashMap;

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
    pub fn get(&mut self, key: String) -> Option<String> {
        self.map.get(&key).cloned()
    }
    pub fn set(&mut self, key: String, value: String) {
        self.map.insert(key, value);
    }
    pub fn remove(&mut self, key: String) {
        self.map.remove(&key);
    }
}
