# Git Account & SSH Key Manager

Git Account & SSH Key Manager is a CLI for developers who need to manage multiple GitHub accounts and SSH keys on one machine.
The `gacc` project helps you switch Git identities safely, configure account-specific SSH access, and avoid authentication mistakes across work and personal repositories.

## What This Tool Does

`gacc` helps you:

- Manage multiple GitHub accounts on one machine
- Use different SSH keys for work and personal repositories
- Switch Git identities per repository
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

Use a dedicated account profile for each identity and activate it inside the repository you are working on.

### Does this tool work for repository-level Git identity switching?

Yes. `gacc` is designed to keep account switching local to the current repository instead of changing your global Git setup.
