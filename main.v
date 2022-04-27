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

		packages << fragments.retrieve_packages(document)
		aur << fragments.retrieve_aur_packages(document)
		symlinks << fragments.retrieve_symlinks(module_name, document)
		directories << fragments.retrieve_directories(document)
		repositories << fragments.retrieve_repositories(document)
		services << fragments.retrieve_services(document)
		scripts << fragments.retrieve_scripts(document)
		groups << fragments.retrieve_groups(document)
		lines << fragments.retrieve_lines(document)
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
