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
    // Given a String, None if not in map, Some(s) if in map
    pub fn get(&mut self, key: String) -> Option<String> {
        // Get wants a reference, and returns Some(&s), cloned clones the reference.
        self.map.get(&key).cloned()
    }
    // Insert, if there exists an old value, replace it
    pub fn set(&mut self, key: String, value: String) {
        // Insert wants two values
        self.map.insert(key, value);
    }
    // Remove a key from the map
    pub fn remove(&mut self, key: String) {
        // Remove wants a reference
        self.map.remove(&key);
    }
}
