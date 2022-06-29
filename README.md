# pj

![pj](bird.png)

Configuration management for localhost

WARNING: only works with Arch Linux for now.

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
