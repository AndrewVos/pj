module fragments

import toml

struct Script {
	name string
}

pub fn retrieve_scripts(document toml.Any) []Script {
	return document.array().map(fn (v toml.Any) Script {
		return Script{
			name: v.value('name').string()
		}
	})
}

pub fn (scripts []Script) execute() {
	println('not supported')
}
