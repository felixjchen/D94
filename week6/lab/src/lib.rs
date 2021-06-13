pub struct KvStore;

impl KvStore {
  pub fn new() -> KvStore {
    KvStore {}
  }

  pub fn get(&self, key: String) -> Option<String> {
    Some("todo".to_string())
  }
  pub fn set(&self, key: String, value: String) {}

  pub fn remove(&self, key: String) {}
}
