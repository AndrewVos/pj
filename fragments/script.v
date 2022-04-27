module fragments

import toml

struct Script {
	name string
}

pub fn retrieve_scripts(document toml.Doc) []Script {
	mut scripts := []Script{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'script' {
			for key, value in top_level_value.as_map() {
				scripts << Script{
					name: value.value('name').string()
				}
			}
		}
	}

	return scripts
}

pub fn (scripts []Script) execute() {
	println('not supported')
}
