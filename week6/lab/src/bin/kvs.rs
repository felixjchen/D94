// extern crate clap;
use clap::{App, Arg, SubCommand};
use std::process;

fn main() {
  let matches = App::new("In Memory KV Store")
    .version("0.1.0")
    .author("Felix C. <felixchen1998@gmail.com>")
    .about("In Memory KV Store")
    .subcommand(
      SubCommand::with_name("set")
        .about("SET K with V")
        .version("0.1.0")
        .author("Felix C. <felixchen1998@gmail.com>")
        .arg(Arg::with_name("key").help("k").required(true).index(1))
        .arg(Arg::with_name("value").help("v").required(true).index(2)),
    )
    .subcommand(
      SubCommand::with_name("get")
        .about("GET K")
        .version("0.1.0")
        .author("Felix C. <felixchen1998@gmail.com>")
        .arg(Arg::with_name("key").help("k").required(true).index(1)),
    )
    .subcommand(
      SubCommand::with_name("rm")
        .about("REMOVE K")
        .version("0.1.0")
        .author("Felix C. <felixchen1998@gmail.com>")
        .arg(Arg::with_name("key").help("k").required(true).index(1)),
    )
    .get_matches();

  if let Some(matches) = matches.subcommand_matches("set") {
    // println!("SET k {}", matches.value_of("key").unwrap());
    // println!("SET v {}", matches.value_of("value").unwrap());
    eprint!("unimplemented");
    process::exit(-1)
  }

  if let Some(matches) = matches.subcommand_matches("get") {
    // println!("GET k: {}", matches.value_of("key").unwrap());
    eprint!("unimplemented");
    process::exit(-1)
  }

  if let Some(matches) = matches.subcommand_matches("rm") {
    // println!("REMOVE k: {}", matches.value_of("key").unwrap());
    eprint!("unimplemented");
    process::exit(-1)
  }

  // No matches, BAD
  process::exit(-1)
}
