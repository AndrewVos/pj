# pj

![pj](bird.png)

> dotfiles configuration management for your machine

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [About](#about)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [Supported objects](#supported-objects)
  - [Pacman packages](#pacman-packages)
  - [Arch User Repository packages](#arch-user-repository-packages)
  - [Brew packages](#brew-packages)
  - [Directory](#directory)
  - [Group](#group)
  - [Script](#script)
  - [Service](#service)
  - [Symlink](#symlink)
- [Shell completions](#shell-completions)
  - [Fish](#fish)
  - [Bash](#bash)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## About

`pj` is a tool to help streamline the configuration of your local machine. It
allows a simple way to:

- Install packages from Homebrew, pacman or AUR
- Symlink files to your home directory (or any directory!)
- Start and stop services
- Run scripts
- Create folders and files

Everything is controlled by simple YAML configuration inside module folders to
keep your concerns separate and easy to manage.

The provided `create-module` CLI makes this step a breeze.

Running `py apply` will sync the configuration to your system.

It works just as well for setting up a new machine as it does managing an
existing one. Just add a new module and run `pj apply` again and the new
configuration will be applied. Existing configuration will be left as is.

## Installation

`pj` is installed with `go` but is provided as a single binary:

```
go install github.com/AndrewVos/pj@latest
```

## Quick start

Let's setup Zsh with the [starship](https://starship.rs/) prompt on a new machine:

```
~/dotfiles
❯ pj create-module zsh
Creating /Users/user/dotfiles/modules...
Creating /Users/user/dotfiles/modules/zsh...
Creating /Users/user/dotfiles/modules/zsh/files...
Creating /Users/user/dotfiles/modules/zsh/configuration.yml...
```

And now add some values to the `configuration.yml`:

```yaml
# install some required packages first
- brew:
    name:
      - zsh
      - starship

# create a `zshrc` file and link it to our home dir
- symlink:
    from: "~/.zshrc"
    to: "zshrc"

# add the shell to `/etc/shells` and change to it
- script:
    command: |
      sudo sh -c 'echo "/usr/local/bin/zsh" >> /etc/shells
      chsh -s /usr/local/bin/zsh'
      echo 'shell changed to zsh!'
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
pj create-module my-module
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
- brew:
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

Create a directory with sudo:

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
