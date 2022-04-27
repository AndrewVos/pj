import os
import toml
import fragments

fn main() {
	module_paths := os.glob('modules/*') or { panic(err) }

	mut packages := []fragments.Package{}
	mut aur := []fragments.Aur{}
	mut symlinks := []fragments.Symlink{}
	mut directories := []fragments.Directory{}
	mut repositories := []fragments.Repository{}
	mut services := []fragments.Service{}
	mut scripts := []fragments.Script{}
	mut groups := []fragments.Group{}
	mut lines := []fragments.Line{}

	for module_path in module_paths {
		module_name := os.base(module_path)

		contents := os.read_file(os.join_path(module_path, 'configuration.toml')) or { panic(err) }
		document := toml.parse_text(contents) or { panic(err) }

		for key, value in document.to_any().as_map() {
			if key == 'package' {
				packages << fragments.retrieve_packages(value)
			} else if key == 'aur' {
				aur << fragments.retrieve_aur_packages(value)
			} else if key == 'symlink' {
				symlinks << fragments.retrieve_symlinks(module_name, value)
			} else if key == 'directory' {
				directories << fragments.retrieve_directories(value)
			} else if key == 'repository' {
				repositories << fragments.retrieve_repositories(value)
			} else if key == 'service' {
				services << fragments.retrieve_services(value)
			} else if key == 'script' {
				scripts << fragments.retrieve_scripts(value)
			} else if key == 'group' {
				groups << fragments.retrieve_groups(value)
			} else if key == 'line' {
				lines << fragments.retrieve_lines(value)
			}
		}
	}

	packages.execute()
	aur.execute()
	symlinks.execute()
	directories.execute()
	repositories.execute()
	services.execute()
	scripts.execute()
	groups.execute()
	lines.execute()
}
