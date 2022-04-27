module fragments

import toml

struct Service {
	name   string
	enable bool
	start  bool
}

pub fn retrieve_services(document toml.Any) []Service {
	return document.array().map(fn (v toml.Any) Service {
		return Service{
			name: v.value('name').string()
			enable: v.value('enable').bool()
			start: v.value('start').bool()
		}
	})
}

pub fn (services []Service) execute() {
	println('not supported')
}
