module fragments

import os
import toml

struct Directory {
	path string
}

fn (d Directory) expand_path() string {
	return os.expand_tilde_to_home(d.path)
}

pub fn (directories []Directory) execute() {
	for directory in directories {
		path := directory.expand_path()

		if !os.is_dir(path) {
			os.mkdir(path) or { panic(err) }
		}
	}
}

pub fn retrieve_directories(document toml.Doc) []Directory {
	mut directories := []Directory{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'directory' {
			for key, value in top_level_value.as_map() {
				directories << Directory{
					path: value.value('path').string()
				}
			}
		}
	}

	return directories
}
