package modules

type Symlink struct {
	Sudo     bool
	From     string
	To       string
	fullFrom string
	fullTo   string
}
