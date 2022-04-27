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

pub fn retrieve_directories(document toml.Any) []Directory {
	return document.array().map(fn (v toml.Any) Directory {
		return Directory{
			path: v.value('path').string()
		}
	})
}
