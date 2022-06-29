# pj

![pj](bird.png)

Configuration management for localhost

WARNING: only works with Arch Linux for now.

## What is this?

When you run ```pj apply``` we sync your configuration to your system,
which means that packages will only get installed once, symlinks will only be created once etc.

## Installation

```
go install github.com/AndrewVos/pj@latest
```

## Configuration

We use the concept of "modules" in `pj`.

A module can be thought of as a group of configuration and files.

Imagine a module called `vim`, with a file structure that looks like:

```
modules/
└── vim
    ├── configuration.yml
        └── files
                └── .vimrc
```

And in `configuration.yml` we have:

```
- pacman:
    name: vim

- symlink:
    from: "~/.vimrc"
    to: ".vimrc"
```

This is just a yml file (and you may notice it's pretty much the same as Ansible).

The first part installs a pacman package called `vim`.

The second part symlinks `~/.vimrc` to the `.vimrc` in your `modules/vim/files/.vimrc`.

Notice that `pj` understands paths relative to your `files` directory inside the module.

Now when you run `pj apply` we will install `vim` and symlink your `.vimrc` to the correct place.

For more examples [take a look at my personal dotfiles repo](https://github.com/AndrewVos/dotfiles).

## Usage

```
pj apply
```

## Supported objects

### Pacman packages

Install a single package:

```yaml
- pacman:
    name: imv
```

Install multiple packages:

```yaml
- pacman:
    name:
      - mpv
      - redis
```

### Arch User Repository packages

Install a single package:

```yaml
- aur:
    name: enpass-bin
```

Install multiple packages:

```yaml
- aur:
    name:
      - enpass-bin
      - google-chrome
      - spotify
      - colorpicker
      - slack-desktop
```

### Brew packages

Install a single package:

```yaml
- brew:
    name: postgresql
```

Install multiple packages:

```yaml
- aur:
    name:
      - postgresql
      - chrome
```

### Directory

Create a directory:

```yaml
- directory:
    path: "~/.my-directory"
```

Create a directory with `sudo`:

```yaml
- directory:
    sudo: true
    path: "/etc/blah"
```

### Group

Add a user to a group:

```yaml
- group:
    user: vos
    name: power
```

### Script

Run a script:

```yaml
- script:
    command: '[[ "$SHELL" = /usr/bin/fish ]] || sudo usermod --shell /usr/bin/fish "$USER"'
```

### Service

Start a service:

```yaml
- service:
    name: "sshd"
    start: true
```

Start and enable a service:

```yaml
- service:
    name: "sshd"
    enable: true
    start: true
```

### Symlink

Symlink a file or directory:

```yaml
- symlink:
    from: "~/.ssh"
    to: ".ssh"
```

Symlink a file or directory with sudo:

```yaml
- symlink:
    sudo: true
    from: "/etc/blah"
    to: "blah"
```

## Shell completions

### Fish

Add the following to your `config.fish`:

```
  pj completion fish | source
```

### Bash

Add the following to your `.bashrc`:

```
source <(pj completion bash)
```
