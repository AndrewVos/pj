module fragments

import toml
import os

struct Symlink {
	module_name string
	become      bool
	from        string
	to          string
}

fn (s Symlink) full_from() string {
	return os.expand_tilde_to_home(s.from)
}

fn (s Symlink) full_to() string {
	if s.to.starts_with('/') {
		return s.to
	}
	return os.join_path(os.getwd(), 'modules', s.module_name, 'files', s.to)
}

pub fn (symlinks []Symlink) execute() {
	mut missing := []Symlink{}

	for symlink in symlinks {
		if os.real_path(symlink.full_from()) != symlink.full_to() {
			missing << symlink
		}
	}

	for symlink in missing {
		if os.is_file(symlink.full_from()) || os.is_link(symlink.full_from()) {
			println('error: file $symlink.full_from() already exists')
			exit(1)
		} else {
			if symlink.become {
				mut process := os.new_process('/bin/sudo')
				process.set_args(['ln', '-s', symlink.full_to(),
					symlink.full_from()])
				process.run()
				process.wait()
				if process.code > 0 {
					panic('error: failed to create symlink $symlink.full_from()')
				}
			} else {
				os.symlink(symlink.full_to(), symlink.full_from()) or { panic(err) }
			}
		}
	}
}

pub fn retrieve_symlinks(module_name string, document toml.Doc) []Symlink {
	mut symlinks := []Symlink{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'symlink' {
			for key, value in top_level_value.as_map() {
				symlinks << Symlink{
					module_name: module_name
					become: value.value('become').default_to(false).bool()
					from: value.value('from').string()
					to: value.value('to').string()
				}
			}
		}
	}

	return symlinks
}
