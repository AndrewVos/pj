module fragments

import toml

struct Line {
	path string
	line string
}

pub fn retrieve_lines(document toml.Any) []Line {
	return document.array().map(fn (v toml.Any) Line {
		return Line{
			path: v.value('path').string()
			line: v.value('line').string()
		}
	})
}

pub fn (lines []Line) execute() {
	println('not supported')
}
