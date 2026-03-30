# How to manage multiple GitHub accounts with SSH keys

Managing multiple GitHub accounts on one machine usually requires separate SSH keys, custom SSH host aliases, and repository-specific Git identity settings.
`gacc` simplifies that process by generating keys, updating SSH config, connecting to GitHub, and switching the active account per repository.

## Why This Is Hard

Using both work and personal GitHub accounts on one computer often causes:

- SSH key conflicts
- Wrong Git author information in commits
- Pushes going to the wrong GitHub account
- Repeated manual edits to `~/.ssh/config`

## Recommended Workflow

### 1. Add each GitHub account

```bash
gacc add work
gacc add personal
```

This creates separate SSH keys and registers them with GitHub.

### 2. Activate the correct account inside each repository

```bash
cd ~/projects/company-repo
gacc activate work

cd ~/projects/side-project
gacc activate personal
```

This updates the current repository so it uses the right SSH route and Git identity.

### 3. Return to your default setup when needed

```bash
gacc deactivate
```

This removes repository-level overrides and restores normal global Git behavior.

## What `gacc` Changes

When you activate an account, `gacc` can:

- Point the repository remote to an account-specific GitHub SSH alias
- Set local `user.name`
- Set local `user.email`

This keeps account switching local to the repository instead of affecting all repositories on the machine.

## When To Use `gacc`

Use `gacc` if you want to:

- Manage multiple GitHub accounts with SSH keys
- Switch between work and personal Git identities safely
- Avoid editing SSH config manually for every new account
- Reduce authentication mistakes across repositories

## Related Questions

### How do I use different SSH keys for different GitHub accounts?

Use a different SSH key and SSH alias for each account.
`gacc` automates the setup and GitHub registration for you.

### How do I avoid Git author conflicts between work and personal repositories?

Apply repository-specific Git settings instead of changing your global Git identity.
`gacc activate [name]` is built for this workflow.

### How do I switch GitHub accounts without breaking existing repositories?

Keep each repository pinned to the right account and only change its local configuration.
That approach avoids accidental global changes.
