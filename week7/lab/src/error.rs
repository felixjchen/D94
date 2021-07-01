use std::ffi::OsString;

#[derive(Debug)]
pub enum KvsError {
  KeyNotFound,
  SerdeError(serde_json::Error),
  IoError(std::io::Error),
  OsStringError(OsString),
  OtherError(String),
}

impl From<OsString> for KvsError {
  fn from(error: OsString) -> Self {
    KvsError::OsStringError(error)
  }
}

impl From<std::io::Error> for KvsError {
  fn from(error: std::io::Error) -> Self {
    KvsError::IoError(error)
  }
}

impl From<serde_json::Error> for KvsError {
  fn from(error: serde_json::Error) -> Self {
    KvsError::SerdeError(error)
  }
}

impl From<String> for KvsError {
  fn from(error: String) -> Self {
    KvsError::OtherError(error)
  }
}

pub type Result<T> = std::result::Result<T, KvsError>;
