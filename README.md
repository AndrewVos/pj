# pj

dotfiles configuration management for your machine

![pj](bird.png)

Your dotfiles are yours - they describe you and are a part of your place on this planet.

They should be beautiful, but they probably aren't right now. You've let them get a bit unorganised, haven't you?

That's ok.

We are here to tell you that your dotfiles can be magnificent once again!

Take our hands and we will show you how to make them shine.

Let's rekindle that love you once had.

## About

`pj` is a tool to help streamline the configuration of your local machine. It
allows you to:

- Install packages from Homebrew, Pacman or AUR
- Symlink files to your home directory (or any directory!)
- Start and stop services
- Run scripts
- Create files and directories

Everything is controlled by YAML configuration which keeps your concerns separate and easy to manage.

## Installation

Binary:

```
wget -O - https://raw.githubusercontent.com/AndrewVos/pj/main/install-from-github | bash
```

With golang:

```
go install github.com/AndrewVos/pj@latest
```

If you're using `fish` shell add the following to your `config.fish`:

```
pj completion fish | source
```

If you're using `bash` add the following to your `.bashrc`:

```
source <(pj completion bash)
```

## Quick start

```
# Create some dotfiles
mkdir dotfiles
cd dotfiles
git init

# Use a task generator to add a symlink from ~/.vimrc to the .vimrc in your dotfiles
pj add symlink module-name --from ~/.vimrc --to .vimrc
touch module/module-name/files/.vimrc

# Apply!
pj apply
```

For some more examples [take a look at my personal dotfiles repo](https://github.com/AndrewVos/dotfiles).

## Task generators

### AUR package

```bash
pj add aur module-name --name package1 --name package2
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- aur:
    name:
    - package1
    - package2
```

### Homebrew package

```bash
pj add brew module-name --name package1 --name package2
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- brew:
    name:
    - package1
    - package2
```

### Pacman package

```bash
pj add pacman module-name --name package1 --name package2
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- pacman:
    name:
    - package1
    - package2
```

### Directory

```bash
pj add directory module-name --path /some/path
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- directory:
    path: /some/path
```

### Directory

```bash
pj add directory module-name --path /some/path --sudo
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- directory:
    path: /some/path
    sudo: true
```

### Group

```bash
pj add group module-name --user some-user --name group-name
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- group:
    name: group-name
    user: some-user
```

### Script

```bash
pj add script module-name --command ls -a
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- script:
    command: ls -a
```

### Service

```bash
pj add service module-name --name service-name --start --enable
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- service:
    enable: true
    name: service-name
    start: true
```

### Symlink

```bash
pj add symlink module-name --from /blah/blah --to blah --sudo
```

```
./modules/module-name/configuration.yml
```

```yml

tasks:
- symlink:
    from: /blah/blah
    sudo: true
    to: blah
```


