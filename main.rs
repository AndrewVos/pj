use serde::{Deserialize, Serialize};
use std::fs;

#[derive(Debug, PartialEq, Serialize, Deserialize)]
struct Pacman {
    name: Vec<String>,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
struct Aur {
    name: Vec<String>,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
struct Fragment {
    pacman: Option<Pacman>,
    aur: Option<Aur>,
}

fn main() -> Result<(), serde_yaml::Error> {
    let contents = fs::read_to_string("../../.dotfiles/fragments/alacritty/configuration.yml")
        .expect("Something went wrong reading the file");

    let fragment: Fragment = serde_yaml::from_str(&contents)?;
    println!("{:?}", fragment);

    Ok(())
}
