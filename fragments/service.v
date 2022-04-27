module fragments

import toml

struct Service {
	name   string
	enable bool
	start  bool
}

pub fn retrieve_services(document toml.Doc) []Service {
	mut services := []Service{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'service' {
			for key, value in top_level_value.as_map() {
				services << Service{
					name: value.value('name').string()
					enable: value.value('enable').bool()
					start: value.value('start').bool()
				}
			}
		}
	}

	return services
}

pub fn (services []Service) execute() {
	println('not supported')
}
