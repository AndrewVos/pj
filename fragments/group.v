module fragments

import toml

struct Group {
	user   string
	groups []string
}

pub fn retrieve_groups(document toml.Doc) []Group {
	mut groups := []Group{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'group' {
			for value in top_level_value.array() {
				groups << Group{
					user: value.value('user').string()
					groups: value.value('groups').array().as_strings()
				}
			}
		}
	}

	return groups
}

pub fn (groups []Group) execute() {
	println('not supported')
}
