module fragments

import os
import toml

struct Repository {
	repository  string
	destination string
}

fn (r Repository) expand_destination() string {
	return os.expand_tilde_to_home(r.destination)
}

pub fn retrieve_repositories(document toml.Any) []Repository {
	return document.array().map(fn (v toml.Any) Repository {
		return Repository{
			repository: v.value('repository').string()
			destination: v.value('destination').string()
		}
	})
}

pub fn (repositories []Repository) execute() {
	for repository in repositories {
		destination := repository.expand_destination()

		if !os.is_dir(destination) {
			mut process := os.new_process('/bin/git')
			process.set_args(['clone', repository.repository, destination])
			process.run()
			process.wait()
			if process.code > 0 {
				panic('error: failed to clone repository $repository.repository')
			}
		}
	}
}
