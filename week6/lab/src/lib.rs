pub struct KvStore;

impl KvStore {
  pub fn new() -> KvStore {
    KvStore {}
  }
  pub fn get(&self, _key: String) -> Option<String> {
    None
  }
  pub fn set(&self, _key: String, _value: String) {}

  pub fn remove(&self, _key: String) {}
}
