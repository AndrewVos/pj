module fragments

import toml

struct Group {
	user   string
	groups []string
}

pub fn retrieve_groups(document toml.Any) []Group {
	return document.array().map(fn (v toml.Any) Group {
		return Group{
			user: v.value('user').string()
			groups: v.value('groups').array().as_strings()
		}
	})
}

pub fn (groups []Group) execute() {
	println('not supported')
}
