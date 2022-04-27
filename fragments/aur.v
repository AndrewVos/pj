module fragments

import os
import toml

struct Aur {
	names []string
}

pub fn retrieve_aur_packages(document toml.Doc) []Aur {
	mut aur_packages := []Aur{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'aur' {
			for aur_package in top_level_value.array() {
				aur_packages << Aur{
					names: aur_package.value('name').array().as_strings()
				}
			}
		}
	}

	return aur_packages
}

fn (packages []Aur) missing() []string {
	mut result := []string{}
	installed := os.execute('pacman -Qq').output.trim('\n').split('\n')

	for package in packages {
		for name in package.names {
			if name !in installed {
				result << name
			}
		}
	}

	return result
}

pub fn (aur_packages []Aur) execute() {
	missing_packages := aur_packages.missing()
	if missing_packages.len > 0 {
		result := os.system('yay -S ${missing_packages.join(' ')}')
		if result != 0 {
			exit(result)
		}
	}
}
