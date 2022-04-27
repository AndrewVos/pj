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

pub fn retrieve_repositories(document toml.Doc) []Repository {
	mut repositories := []Repository{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'repository' {
			for repository in top_level_value.array() {
				repositories << Repository{
					repository: repository.value('repository').string()
					destination: repository.value('destination').string()
				}
			}
		}
	}

	return repositories
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
