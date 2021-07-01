use clap::{App, Arg, SubCommand};
use kvs::{KvStore, Result};
use std::process;

fn main() -> Result<()> {
  let matches = App::new(env!("CARGO_PKG_NAME"))
    .version(env!("CARGO_PKG_VERSION"))
    .author(env!("CARGO_PKG_AUTHORS"))
    .about(env!("CARGO_PKG_DESCRIPTION"))
    .subcommand(
      SubCommand::with_name("set")
        .about("Set K with V")
        .version(env!("CARGO_PKG_VERSION"))
        .author(env!("CARGO_PKG_AUTHORS"))
        .arg(Arg::with_name("key").help("k").required(true).index(1))
        .arg(Arg::with_name("value").help("v").required(true).index(2)),
    )
    .subcommand(
      SubCommand::with_name("get")
        .about("Get K")
        .version(env!("CARGO_PKG_VERSION"))
        .author(env!("CARGO_PKG_AUTHORS"))
        .arg(Arg::with_name("key").help("k").required(true).index(1)),
    )
    .subcommand(
      SubCommand::with_name("rm")
        .about("Remove K")
        .version(env!("CARGO_PKG_VERSION"))
        .author(env!("CARGO_PKG_AUTHORS"))
        .arg(Arg::with_name("key").help("k").required(true).index(1)),
    )
    .get_matches();

  if let Some(matches) = matches.subcommand_matches("set") {
    let key = matches.value_of("key").unwrap().to_string();
    let value = matches.value_of("value").unwrap().to_string();

    let mut kvs = KvStore::open("")?;
    kvs.set(key, value)?;
    process::exit(0)
  }

  if let Some(matches) = matches.subcommand_matches("get") {
    let key = matches.value_of("key").unwrap().to_string();

    let mut kvs = KvStore::open("")?;
    match kvs.get(key)? {
      None => println!("Key not found"),
      Some(value) => println!("{}", value),
    };

    process::exit(0)
  }

  if let Some(matches) = matches.subcommand_matches("rm") {
    let key = matches.value_of("key").unwrap().to_string();

    let mut kvs = KvStore::open("")?;
    kvs.remove(key)?;
    process::exit(0)
  }

  // No matches, BAD
  eprint!("unimplemented");
  process::exit(-1)
}
