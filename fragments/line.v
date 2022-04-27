module fragments

import toml

struct Line {
	path string
	line string
}

pub fn retrieve_lines(document toml.Doc) []Line {
	mut lines := []Line{}

	for top_level_key, top_level_value in document.to_any().as_map() {
		if top_level_key == 'line' {
			for line in top_level_value.array() {
				lines << Line{
					path: line.value('path').string()
					line: line.value('line').string()
				}
			}
		}
	}

	return lines
}

pub fn (lines []Line) execute() {
}
