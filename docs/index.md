# Git Account & SSH Key Manager

Git Account & SSH Key Manager is a CLI for developers who need to manage multiple GitHub accounts and SSH keys on one machine.
The `gacc` project helps you switch Git identities safely, configure account-specific SSH access, and avoid authentication mistakes across work and personal repositories.

## What This Tool Does

`gacc` helps you:

- Manage multiple GitHub accounts on one machine
- Use different SSH keys for work and personal repositories
- Switch Git identities per repository
- Automatically apply Git identities by directory
- Set a global default Git identity for the machine
- Avoid pushing with the wrong GitHub account

## Start Here

- [How to manage multiple GitHub accounts with SSH keys](how-to-manage-multiple-github-accounts-with-ssh-keys.md)
- [GitHub Releases](https://github.com/soongeo/gacc/releases)
- [Project README](../README.md)

## Common Questions

### How do I manage multiple GitHub accounts with SSH keys?

Create separate SSH identities for each account and apply the correct one to each repository.
`gacc` automates that workflow.

### How do I switch between work and personal GitHub accounts?

Use a dedicated account profile for each identity and choose one of three workflows:

- `gacc activate [name]` for per-repository manual switching
- `gacc auto add [name] [directory]` for directory-based automatic switching
- `gacc global activate [name]` for a machine-wide default

Run `gacc status` to see which layer is active for the current directory.

### Does this tool work for repository-level Git identity switching?

Yes. `gacc` is designed to keep account switching local to the current repository instead of changing your global Git setup.

### Can this tool automatically switch accounts by folder?

Yes. `gacc auto add` uses Git `includeIf` rules so directories like `~/Work` and `~/Personal` can automatically use different identities.

### What is the priority between local, automatic, and global settings?

`gacc` applies the most specific rule first:

1. Local repository settings from `gacc activate`
2. Directory-based automatic settings from `gacc auto add`
3. Global defaults from `gacc global activate`
4. Any other plain Git defaults
