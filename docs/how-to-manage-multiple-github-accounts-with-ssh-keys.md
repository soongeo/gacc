# How to manage multiple GitHub accounts with SSH keys

Managing multiple GitHub accounts on one machine usually requires separate SSH keys, custom SSH host aliases, and repository-specific Git identity settings.
`gacc` simplifies that process by generating keys, updating SSH config, connecting to GitHub, and switching the active account locally, automatically by directory, or globally for the whole machine.

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

### 2. Choose how the account should be applied

#### Option A: Per-repository manual switching

```bash
cd ~/projects/company-repo
gacc activate work

cd ~/projects/side-project
gacc activate personal
```

This updates the current repository so it uses the right SSH route and Git identity.

#### Option B: Directory-based automatic switching

```bash
gacc auto add work ~/Work
gacc auto add personal ~/Personal
```

This uses Git `includeIf` so repositories under those folders automatically use the configured account.

#### Option C: Global default account

```bash
gacc global activate personal
```

This sets a default identity for Git operations everywhere unless a more specific local or automatic rule overrides it.

### 3. Check which account layer is active

```bash
gacc status
```

This shows the local, automatic, and global account layers for the current directory.

### 4. Return to your default setup when needed

```bash
gacc deactivate
gacc global deactivate
```

Use `gacc deactivate` to clear repository-level overrides and `gacc global deactivate` to clear the machine-wide default.

## What `gacc` Changes

When you manually activate an account with `gacc activate`, `gacc` can:

- Point the repository remote to an account-specific GitHub SSH alias
- Set local `user.name`
- Set local `user.email`

When you use `gacc auto add`, `gacc` creates Git `includeIf` rules that apply account settings by directory.

When you use `gacc global activate`, `gacc` sets global `user.name`, `user.email`, and `core.sshCommand`.

## Priority Rules

`gacc` applies the most specific layer first:

1. Local repository settings from `gacc activate`
2. Directory-based automatic settings from `gacc auto add`
3. Global defaults from `gacc global activate`
4. Any other plain Git defaults

In practice:

- Local settings win over automatic and global settings
- Automatic directory rules win over global defaults
- Global settings act as the fallback

Run `gacc status` to inspect the active layers for the current directory.

## When To Use `gacc`

Use `gacc` if you want to:

- Manage multiple GitHub accounts with SSH keys
- Switch between work and personal Git identities safely
- Automatically apply accounts by folder
- Keep one default identity for everything else
- Avoid editing SSH config manually for every new account
- Reduce authentication mistakes across repositories

## Related Questions

### How do I use different SSH keys for different GitHub accounts?

Use a different SSH key and SSH alias for each account.
`gacc` automates the setup and GitHub registration for you.

### How do I avoid Git author conflicts between work and personal repositories?

Apply repository-specific Git settings instead of changing your global Git identity.
`gacc activate [name]` is built for this workflow.

### How do I automatically use my work account inside `~/Work`?

Run `gacc auto add work ~/Work`.
That directory and its repositories will automatically use the configured account settings.

### How do I switch GitHub accounts without breaking existing repositories?

Keep each repository pinned to the right account and only change its local configuration.
That approach avoids accidental global changes.
