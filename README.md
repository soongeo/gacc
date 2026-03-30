# Git Account & SSH Key Manager

`gacc` is the project name for this Git Account & SSH Key Manager.

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
- Deactivate repository-specific overrides and return to global defaults
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

### `gacc deactivate`

Removes repository-specific overrides so the repository goes back to your default Git behavior.

```bash
cd client-project
gacc deactivate
```

### `gacc delete [name]`

Deletes the local SSH key and config for an account and attempts to remove the registered public key from GitHub.

```bash
gacc delete work
```

## Common Use Cases

### Use different GitHub accounts for work and personal repositories

Create one account profile for `work` and another for `personal`, then activate the right one inside each repository.

### Avoid SSH conflicts on one machine

`gacc` creates account-specific SSH aliases so multiple GitHub accounts can coexist without editing SSH settings by hand every time.

### Keep repository identities local

`gacc activate` writes repository-level settings instead of changing your global Git identity for every project.

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

### Does `gacc` change my global Git configuration?

No. `gacc activate` is designed to apply account settings to the current repository so your global defaults stay intact.

### Can `gacc` help prevent pushing with the wrong GitHub account?

Yes. It updates the Git remote host alias and repository-level Git identity so the selected repository uses the intended account.

## Troubleshooting

### GitHub authentication does not start

Check that your network allows GitHub device flow authentication and try running `gacc add [name]` again.

### The current directory is not recognized as a Git repository

Run `git status` first and make sure you are inside a repository before using `gacc activate` or `gacc deactivate`.

### The wrong account still appears active

Run `gacc list` to inspect registered accounts, then verify that the repository `origin` remote uses the expected GitHub host alias.

## Documentation

- [Git Account & SSH Key Manager Docs](docs/index.md)
- [How to manage multiple GitHub accounts with SSH keys](docs/how-to-manage-multiple-github-accounts-with-ssh-keys.md)

## Releases

Download binaries and release builds from the [GitHub Releases](https://github.com/soongeo/gacc/releases) page.

## Contributing

Questions, bug reports, and pull requests are welcome at the [GitHub repository](https://github.com/soongeo/gacc).
