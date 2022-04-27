module fragments

import os
import toml

struct Package {
	names []string
}

pub fn retrieve_packages(document toml.Any) []Package {
	return document.array().map(fn (v toml.Any) Package {
		return Package{
			names: v.value('name').array().as_strings()
		}
	})
}

fn (packages []Package) missing() []string {
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

pub fn (packages []Package) execute() {
	missing_packages := packages.missing()
	if missing_packages.len > 0 {
		result := os.system('sudo pacman -S ${missing_packages.join(' ')}')
		if result != 0 {
			exit(result)
		}
	}
}
