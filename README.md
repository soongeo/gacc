# Git Account & SSH Key Manager

[![Build Status](https://github.com/soongeo/gacc/actions/workflows/ci.yml/badge.svg)](https://github.com/soongeo/gacc/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/soongeo/gacc)](https://goreportcard.com/report/github.com/soongeo/gacc) [![Go Reference](https://pkg.go.dev/badge/github.com/soongeo/gacc.svg)](https://pkg.go.dev/github.com/soongeo/gacc) [![Go Version](https://img.shields.io/github/go-mod/go-version/soongeo/gacc)](https://github.com/soongeo/gacc) [![Latest Release](https://img.shields.io/github/v/release/soongeo/gacc)](https://github.com/soongeo/gacc/releases) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gacc` is the project name for this Git Account & SSH Key Manager.

# gacc

`gacc` is a CLI for managing multiple GitHub accounts and SSH keys on one machine.
It helps developers switch between personal and work Git identities safely, configure SSH access per account, and avoid pushing or committing with the wrong credentials.

## What Is gacc?

`gacc` is a Git account manager for developers who use more than one GitHub account on the same computer.
It automates SSH key generation, GitHub authentication, SSH config updates, and per-repository Git identity switching.

## Who Is This For?

`gacc` is useful for developers who need to:

- Manage multiple GitHub accounts on one machine
- Use separate SSH keys for work and personal repositories
- Switch Git identities per repository without changing global Git settings
- Reduce GitHub authentication mistakes and SSH host conflicts

## What Problem Does It Solve?

Managing multiple GitHub accounts usually means editing `~/.ssh/config`, creating and naming SSH keys manually, and remembering which repositories should use which identity.
`gacc` turns that into a repeatable workflow so you can:

- Manage multiple GitHub accounts without overwriting SSH settings
- Switch Git identities safely between repositories
- Route Git remotes through account-specific SSH aliases
- Keep local repository settings isolated from your global Git config

## Features

- Generate and register SSH keys for each Git account
- Authenticate with GitHub using the device flow
- Update `~/.ssh/config` with account-specific GitHub host aliases
- Activate a selected Git identity for the current Git repository
- Activate a selected Git identity globally for all Git operations
- Automatically apply accounts by directory using Git `includeIf`
- Deactivate repository-specific overrides and return to global defaults
- Show local, auto, and global account status for the current directory
- Display, rename, cache, backup, and restore managed account keys
- Delete local account records and remove uploaded SSH keys from GitHub

## Installation

### macOS and Linux with Homebrew

```bash
brew tap soongeo/gacc https://github.com/soongeo/gacc
brew install gacc
```

### Install with Go

```bash
go install github.com/soongeo/gacc@latest
```

### Manual Installation

Prebuilt binaries for macOS, Linux, and Windows are available on the [Releases](https://github.com/soongeo/gacc/releases) page.

## Quick Start

```bash
gacc add work
cd your-repository
gacc activate work
```

This flow creates an SSH key for the `work` account, uploads it to GitHub, and applies that identity to the current repository.

For a directory-based workflow, you can configure automatic switching:

```bash
gacc add work
gacc add personal
gacc auto add work ~/Work
gacc auto add personal ~/Personal
```

For a machine-wide default identity, you can configure a global account:

```bash
gacc global activate work
```

## Commands

### `gacc add [name]`

Adds a new Git account, generates an SSH key, starts GitHub device authentication, uploads the public key, and stores profile data locally.

```bash
gacc add work
```

### `gacc list`

Lists the Git accounts registered on your machine and highlights the account currently active for the repository when possible.

```bash
gacc list
```

### `gacc activate [name]`

Applies a selected account to the current Git repository by updating the `origin` remote and local `user.name` and `user.email` values.

```bash
cd client-project
gacc activate work
```

If you omit `[name]`, `gacc` will try to use the currently active account or prompt you to choose one.

### `gacc deactivate`

Removes repository-specific overrides so the repository goes back to your default Git behavior.

```bash
cd client-project
gacc deactivate
```

### `gacc global activate [name]`

Applies a selected account globally by setting global `user.name`, `user.email`, and `core.sshCommand`.

```bash
gacc global activate work
```

### `gacc global deactivate`

Clears global `gacc`-managed Git identity and SSH command settings.

```bash
gacc global deactivate
```

### `gacc auto add [name] [directory]`

Automatically applies an account for repositories under a directory using Git `includeIf`.

```bash
gacc auto add work ~/Work
```

### `gacc auto list`

Lists configured automatic directory-based account rules.

```bash
gacc auto list
```

### `gacc auto remove [name] [directory]`

Removes an automatic directory-based account rule.

```bash
gacc auto remove work ~/Work
```

### `gacc status`

Shows the local, automatic, and global account state for the current directory.

```bash
gacc status
```

### `gacc display [name]`

Displays the public SSH key for an account.

```bash
gacc display work
```

### `gacc rename [old-name] [new-name]`

Renames an account alias and updates associated key files and config references.

```bash
gacc rename work company
```

### `gacc cache [name]`

Adds an account SSH key to `ssh-agent`.

```bash
gacc cache work
```

### `gacc backup [archive]`

Backs up `gacc` config and managed SSH keys into a `tar.gz` archive.

```bash
gacc backup
```

### `gacc restore [archive]`

Restores `gacc` config and managed SSH keys from a backup archive.

```bash
gacc restore ./gacc-backup-20260330153000.tar.gz
```

### `gacc delete [name]`

Deletes the local SSH key and config for an account and attempts to remove the registered public key from GitHub.

```bash
gacc delete work
```

## Manual Vs Auto Vs Global

`gacc` now supports three ways to apply an identity:

- `local`: Explicitly set an account for the current repository with `gacc activate`
- `auto`: Automatically apply an account by directory with `gacc auto add`
- `global`: Set a machine-wide default with `gacc global activate`

### When To Use Each

- Use `local` when you want to choose the account per repository.
- Use `auto` when everything under a folder such as `~/Work` should always use the same account.
- Use `global` when you want a default account everywhere unless a more specific rule overrides it.

### Priority Rules

`gacc` resolves settings from most specific to least specific:

1. Local repository settings from `gacc activate`
2. Directory-based automatic settings from `gacc auto add`
3. Global defaults from `gacc global activate`
4. Any other plain Git defaults outside `gacc`

In practice, that means:

- Local `user.name` and `user.email` override auto and global values
- Auto rules override global defaults for matching directories
- Global settings act as the fallback when no local or auto rule applies

Use `gacc status` to see which layer is active for the current directory.

## Examples

### Per-Repository Manual Switching

```bash
gacc add work
gacc add personal

cd ~/src/client-a
gacc activate work

cd ~/src/oss-project
gacc activate personal
```

### Directory-Based Automatic Switching

```bash
gacc add work
gacc add personal

gacc auto add work ~/Work
gacc auto add personal ~/Personal

cd ~/Work/client-a
gacc status
```

### Global Default With Local Override

```bash
gacc global activate personal

cd ~/Work/client-a
gacc activate work
gacc status
```

## Common Use Cases

### Use different GitHub accounts for work and personal repositories

Create one account profile for `work` and another for `personal`, then activate the right one inside each repository.

### Avoid SSH conflicts on one machine

`gacc` creates account-specific SSH aliases so multiple GitHub accounts can coexist without editing SSH settings by hand every time.

### Keep repository identities local

`gacc activate` writes repository-level settings instead of changing your global Git identity for every project.

### Automatically apply work and personal identities by folder

Use `gacc auto add` to map directories such as `~/Work` and `~/Personal` to different accounts.

### Set one default account for everything

Use `gacc global activate` when one identity should be the default unless a repository or directory rule overrides it.

## Comparison

### Manual setup

Manual setup usually requires:

- Generating SSH keys manually
- Editing `~/.ssh/config` by hand
- Updating Git remotes for host aliases
- Remembering to change local `user.name` and `user.email`

### With `gacc`

`gacc` combines those steps into a single CLI workflow designed for developers managing multiple GitHub accounts and SSH keys.

## FAQ

### How do I manage multiple GitHub accounts on one computer?

Use separate SSH keys and SSH host aliases for each account, then apply the correct Git identity per repository.
`gacc` automates this setup and switching workflow.

### How do I use different SSH keys for work and personal GitHub accounts?

Create one `gacc` account per identity, such as `work` and `personal`.
Each account gets its own SSH key and GitHub SSH alias.

### How do I switch Git identities between repositories?

Change into the repository directory and run `gacc activate [name]`.
This updates the repository to use the selected account.

### How do I automatically use my work account inside `~/Work`?

Run `gacc auto add work ~/Work`.
Then repositories under that directory will automatically use the configured account settings.

### Does `gacc` change my global Git configuration?

`gacc activate` does not change your global Git configuration.
If you want a global default identity, use `gacc global activate [name]`.

### Can `gacc` help prevent pushing with the wrong GitHub account?

Yes. It updates the Git remote host alias and repository-level Git identity so the selected repository uses the intended account.

## Troubleshooting

### GitHub authentication does not start

Check that your network allows GitHub device flow authentication and try running `gacc add [name]` again.

### The current directory is not recognized as a Git repository

Run `git status` first and make sure you are inside a repository before using `gacc activate` or `gacc deactivate`.

### The wrong account still appears active

Run `gacc status` to inspect the current local, automatic, and global account layers.
You can also run `gacc list` to see which registered accounts are marked as `local`, `auto`, or `global`.

## Documentation

- [Git Account & SSH Key Manager Docs](docs/index.md)
- [How to manage multiple GitHub accounts with SSH keys](docs/how-to-manage-multiple-github-accounts-with-ssh-keys.md)

## Releases

Download binaries and release builds from the [GitHub Releases](https://github.com/soongeo/gacc/releases) page.

## Contributing

Questions, bug reports, and pull requests are welcome at the [GitHub repository](https://github.com/soongeo/gacc).
