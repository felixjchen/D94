use clap::{App, Arg, SubCommand};
use kvs::KvStore;
use std::process;

fn main() {
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
    println!(
      "SET k {} v {}",
      matches.value_of("key").unwrap(),
      matches.value_of("value").unwrap()
    );
    eprint!("unimplemented");
    process::exit(-1)
  }

  if let Some(matches) = matches.subcommand_matches("get") {
    println!("GET k: {}", matches.value_of("key").unwrap());
    eprint!("unimplemented");
    process::exit(-1)
  }

  if let Some(matches) = matches.subcommand_matches("rm") {
    println!("REMOVE k: {}", matches.value_of("key").unwrap());
    eprint!("unimplemented");
    process::exit(-1)
  }

  // No matches, BAD
  eprint!("unimplemented");
  process::exit(-1)
}
