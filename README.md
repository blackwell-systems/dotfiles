# Dotfiles & Vault Setup

This repository contains my personal dotfiles for **macOS** and **Lima** (Linux), used to configure my development environment consistently across both platforms. The dotfiles include configurations for **Zsh**, **Powerlevel10k**, **Homebrew**, **Claude helpers**, and a **Bitwarden-based vault bootstrap** for SSH keys, AWS config/credentials, and environment secrets.

---

## Directory Structure

The dotfiles are organized as follows:

```text
~/workspace/dotfiles
├── bootstrap-dotfiles.sh     # Shared symlink bootstrap (zshrc, p10k, Ghostty, Claude)
├── bootstrap-lima.sh         # Lima / Linux-specific bootstrap wrapper
├── bootstrap-mac.sh          # macOS-specific bootstrap wrapper
├── Brewfile                  # Unified Homebrew bundle (macOS + Lima)
├── ghostty
│   └── config                # Ghostty terminal config
├── lima
│   └── lima.yaml             # Lima VM config (host-side)
├── vault
│   ├── bootstrap-vault.sh    # Orchestrates all Bitwarden restores
│   ├── restore-ssh.sh        # Restores SSH keys from Bitwarden
│   ├── restore-aws.sh        # Restores ~/.aws/config & ~/.aws/credentials
│   └── restore-env.sh        # Restores environment secrets to ~/.local
└── zsh
    ├── p10k.zsh              # Powerlevel10k theme config
    └── zshrc                 # Main Zsh configuration
```

Key pieces:

- **zsh/zshrc**: Main Zsh configuration file  
- **zsh/p10k.zsh**: Powerlevel10k theme configuration  
- **ghostty/config**: Ghostty terminal configuration  
- **vault/**: Bitwarden-based secure bootstrap for SSH, AWS, and environment secrets  
- **Brewfile**: Shared Homebrew definition used by both macOS and Lima bootstrap scripts  
- **Claude Workspace Symlink** inside `bootstrap-dotfiles.sh` ensures that both macOS and Lima point to the shared workspace directory:

  ```
  ~/.claude → ~/workspace/.claude
  ```

---

## Global Prerequisites

On **both macOS and Lima/Linux**, you’ll eventually want:

- **Zsh** as your login shell  
- **Homebrew** (macOS or Linuxbrew)  
- **Git**  
- **Bitwarden CLI** (`bw`)  
- **jq** (for JSON manipulation)  
- **AWS CLI v2** (for AWS workflows)  
- **Claude Code (CLI)** via Homebrew (installed from Brewfile on both macOS + Lima)

You can install most of these via Homebrew (after the basic bootstrap is done).

---

## Bootstrap Overview

There are two big pillars:

1. **Dotfiles / Shell bootstrap**

   Handled by:

   - `bootstrap-dotfiles.sh`
   - `bootstrap-mac.sh`
   - `bootstrap-lima.sh`

   Goal: consistent Zsh + p10k + plugins + Ghostty config + Claude across host and Lima.

2. **Vault / Secure secrets bootstrap (Bitwarden)**

   Handled by:

   - `vault/bootstrap-vault.sh`
   - `vault/restore-ssh.sh`
   - `vault/restore-aws.sh`
   - `vault/restore-env.sh`

   Goal: restore **SSH keys**, **AWS config/credentials**, and **env secrets** from Bitwarden.

---

## Bootstrapping macOS from Scratch

1. **Create workspace directory**

```bash
mkdir -p ~/workspace
cd ~/workspace
```

2. **Clone dotfiles repo**

```bash
git clone git@github.com:your-username/dotfiles.git
cd ~/workspace/dotfiles
```

3. **Run macOS bootstrap**

```bash
./bootstrap-mac.sh
```

Typical responsibilities of `bootstrap-mac.sh`:

- Install **Xcode Command Line Tools** (if missing).  
- Install or update **Homebrew**.  
- Ensure `brew` is on `PATH`.  
- Run the **shared Brewfile**:

  ```bash
  brew bundle --file="$DOTFILES_DIR/Brewfile"
  ```

- Run `bootstrap-dotfiles.sh` to create symlinks:

  - `~/.zshrc    → ~/workspace/dotfiles/zsh/zshrc`  
  - `~/.p10k.zsh → ~/workspace/dotfiles/zsh/p10k.zsh`  
  - `~/.claude   → ~/workspace/.claude`  
  - Ghostty config symlink into:

    ```
    ~/Library/Application Support/com.mitchellh.ghostty/config
    ```

4. **Open a new terminal**

This ensures the new `~/.zshrc`, Powerlevel10k, plugins, and Claude integration are picked up.

---

## Bootstrapping Lima / Linux Guest

Assuming your Lima VM shares `~/workspace` from macOS:

1. **Start Lima**

```bash
limactl start ~/workspace/dotfiles/lima/lima.yaml
limactl shell lima-dev-ubuntu
```

2. **Run the Lima bootstrap**

```bash
cd ~/workspace/dotfiles
./bootstrap-lima.sh
```

Typical responsibilities of `bootstrap-lima.sh`:

- Install essential packages (`git`, `zsh`, etc.).  
- Install **Linuxbrew** if missing.  
- Ensure `brew` is on `PATH`.  
- Run the **same Brewfile** used by macOS.  
- Call `bootstrap-dotfiles.sh` to symlink everything, including:

  - `~/.claude → ~/workspace/.claude`
  - shared Zsh files
  - shared p10k config

3. **Restart Lima shell**

---

## Dotfiles Bootstrap Details

### `bootstrap-dotfiles.sh`

This file creates all unified symlinks:

- `~/.zshrc    -> $DOTFILES_DIR/zsh/zshrc`  
- `~/.p10k.zsh -> $DOTFILES_DIR/zsh/p10k.zsh`  
- Ghostty config (macOS only)  
- **Claude workspace symlink**:

  ```
  if [ ! -L "$HOME/.claude" ]; then
      ln -s ~/workspace/.claude ~/.claude
  fi
  ```

This ensures Claude CLI sees a *shared* workspace on both platforms.

---

## Homebrew & Brewfile

The **Brewfile** is shared by both macOS and Lima and includes:

- Core CLI tools  
- Zsh + plugins  
- Docker  
- Lima  
- jq, awscli, bitwarden-cli  
- Claude Code for both macOS + Linux  
- macOS-only casks (ignored automatically on Linux)

Example (abbreviated):

```ruby
brew "git"
brew "zsh"
brew "tmux"
brew "zellij"
brew "node"
brew "docker"
brew "lima"
brew "powerlevel10k"
brew "zsh-autosuggestions"
brew "zsh-syntax-highlighting"
brew "jq"
brew "awscli"
brew "bitwarden-cli"
brew "claude-code"

on_macos do
  cask "ghostty"
  cask "claude-code"
  cask "font-meslo-for-powerlevel10k"
  cask "microsoft-edge"
  cask "nosql-workbench"
  cask "mongodb-compass"
  cask "rectangle"
  cask "vscodium"
end
```

You can regenerate:

```bash
brew bundle dump --force --file=./Brewfile
```

and prune as needed.

---

## Vault / Bitwarden Bootstrap

Lives under:

```
~/workspace/dotfiles/vault
```

Restores:

- **SSH keys**  
- **AWS config & credentials**  
- **Environment secrets**  

via Bitwarden Secure Notes.

Same flow on macOS and Lima.

---

## Pushing Files into Bitwarden (One-Time)

Scripts show exactly how to push:

- `.aws/config`  
- `.aws/credentials`  
- SSH keys  
- environment secrets  

into appropriate Secure Notes:

- `SSH-GitHub-Enterprise`  
- `SSH-GitHub-Blackwell`  
- `AWS-Config`  
- `AWS-Credentials`  
- `Environment-Secrets`

---

## Using the Dotfiles Day-to-Day

### Aliases (defined in zshrc)

- `gst` → `git status`  
- `gco` → `git checkout`  
- `gp`  → `git push`  
- `gl`  → `git pull`  
- `cws` → `cd ~/workspace`  
- `ccode` → `cd ~/workspace/code`  
- `cwhite` → `cd ~/whitepapers`  
- `cpat` → `cd ~/patent-pool`

### Claude helpers

- `claude-bedrock "prompt"`  
- `claude-max "prompt"`  
- `claude-run bedrock "prompt"`  
- `claude-run max "prompt"`

### Environment secrets

```bash
source ~/.local/load-env.sh
```

---

## Troubleshooting

### Claude workspace not detected

Make sure:

```
ls -ld ~/.claude
# should be → ~/workspace/.claude
```

If not, rerun:

```bash
~/workspace/dotfiles/bootstrap-dotfiles.sh
```

### Bitwarden issues

You can always reset:

```bash
bw logout
bw login
export BW_SESSION="$(bw unlock --raw)"
bw sync --session "$BW_SESSION"
```

---

## License

This repository is licensed under the **MIT License**.

By following this guide, you can fully restore your **dotfiles**, **SSH keys**, **AWS configuration**, **packages via Brewfile**, **Claude workspace**, and **environment secrets** across macOS and Lima/Linux in a reproducible, vault-backed, fully unified way.
