# 🐙 gacc - Git Account & SSH Key Manager

`gacc` is a powerful CLI tool designed to simplify managing multiple Git accounts (and their SSH keys) on a single machine. It seamlessly handles SSH key generation, GitHub authentication, and local repository context switching so you never accidentally commit with the wrong email or push with the wrong SSH key again.

## ✨ Features

- **Effortless Key Management:** Automatically generates SSH keys and configures `~/.ssh/config` for each distinct Git identity.
- **GitHub Device Flow:** Easily authenticate and add SSH keys directly to your GitHub account from the CLI without manual web navigation.
- **Seamless Context Switching:** Apply specific Git accounts (name, email, and localized SSH commands) to individual local repositories, without polluting global Git settings.
- **Full Lifecycle Control:** Add, list, activate, deactivate, and delete Git identities safely and consistently.

## 📦 Installation

### macOS & Linux (Homebrew)

The easiest way to install `gacc` is via Homebrew. 

```bash
brew tap soongeo/gacc https://github.com/soongeo/gacc
brew install gacc
```

### Go Developers

If you already have a Go environment (`>= 1.22`) set up, you can compile and install it directly via `go install`:

```bash
go install github.com/soongeo/gacc@latest
```

### Manual Installation

You can also download the pre-compiled binaries for macOS, Linux, and Windows from the [Releases](https://github.com/soongeo/gacc/releases) page.

---

## 🚀 Usage

`gacc` provides an intuitive set of commands for day-to-day use. Run `gacc help` to see all available commands.

### `gacc add`
Add a new Git account and generate its corresponding SSH key. The command will launch a GitHub Device flow to authorize and automatically upload the SSH key to your account.
```bash
gacc add
```

### `gacc list`
List all the Git accounts currently registered and managed by `gacc` on your system.
```bash
gacc list
```

### `gacc activate`
Activate a specific `gacc` account for the **current Git repository**. It intelligently overrides the local `.git/config` (`user.name`, `user.email`, `core.sshCommand`), guaranteeing that commits and pushes use the designated identity.
```bash
cd custom-client-project
gacc activate
```

### `gacc deactivate`
Removes the `gacc` local configuration from the current Git repository, reverting it to your default global Git behavior.
```bash
cd custom-client-project
gacc deactivate
```

### `gacc delete`
Deletes a registered Git account and its SSH key locally. It will also clean up your `~/.ssh/config` and attempt to remove the specified key from your GitHub account.
```bash
gacc delete
```

## 🛠 Contributing

Questions, bug reports, and pull requests are welcome! 
Feel free to open an issue or submit a PR on our [GitHub repository](https://github.com/soongeo/gacc).

## 📄 License

This project is open-source and available under the terms of the [MIT License](LICENSE).
