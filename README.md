# Dotfiles Setup

This repository contains my personal dotfiles for *macOS* and *Lima* (Linux), used to configure my development environment consistently across both platforms. The dotfiles include configurations for *Zsh, **Powerlevel10k, **Homebrew, **Claude helpers*, and more.

## Directory Structure

The dotfiles are organized as follows:

text
~/workspace/dotfiles
├── bootstrap-dotfiles.sh
├── bootstrap-lima.sh
├── bootstrap-mac.sh
├── Brewfile
├── ghostty
│   └── config
├── lima
│   └── lima.yaml
├── README.md
└── zsh
    ├── p10k.zsh
    └── zshrc


- *zshrc*: Main Zsh configuration file  
- *p10k.zsh*: Powerlevel10k theme configuration  
- *ghostty/config*: Ghostty terminal configuration  

Symlinks in your home directory (~/.zshrc, ~/.p10k.zsh, etc.) will point to these files.

## Prerequisites

1. *macOS* or *Linux (Lima)*
2. *Homebrew* (Linuxbrew on Lima)
3. *Zsh* set as the default shell

## Setup Instructions

### 1. Clone the Repository

Clone into ~/workspace/dotfiles:

bash
cd ~/workspace
git clone https://github.com/your-username/dotfiles.git


### 2. Create Symlinks

All dotfiles exist in their non-dot form inside the repo (e.g., zshrc instead of .zshrc).  
The bootstrap script creates the expected home-directory dotfiles.

Run:

bash
cd ~/workspace/dotfiles
./bootstrap-dotfiles.sh


The script will:

- Create symlinks in your home directory:
  - ~/.zshrc → ~/workspace/dotfiles/zsh/zshrc
  - ~/.p10k.zsh → ~/workspace/dotfiles/zsh/p10k.zsh
- Ensure *Claude helpers* and *Zsh plugins* are wired in

### 3. Verify Symlinks

bash
ls -lash ~/.zshrc
ls -lash ~/.p10k.zsh


You should see symlinks pointing to the files inside the dotfiles repository.

Reload Zsh:

bash
source ~/.zshrc


### 4. Install Dependencies (Optional)

#### macOS

bash
brew install powerlevel10k
brew install zsh-autosuggestions
brew install zsh-syntax-highlighting


#### Lima (Linux)

bash
# Install Homebrew (Linuxbrew) if it's not installed already
if ! command -v brew >/dev/null 2>&1; then
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Install Powerlevel10k and Zsh plugins
brew install powerlevel10k zsh-autosuggestions zsh-syntax-highlighting


### 5. Confirm Appearance

You should now have:

- *Syntax highlighting*
- *Autosuggestions*
- *Powerlevel10k prompt*

---

## Key Configuration Files

- *~/.zshrc* — Main configuration for Zsh  
- *~/.p10k.zsh* — Powerlevel10k configuration  
- *Claude helpers* — Functions for controlling which Claude mode runs

---

## Using the Dotfiles

### Switching Claude Modes

Available functions:

- *claude-bedrock* — Uses AWS Bedrock  
- *claude-max* — Uses Anthropic Max  
- *claude-run MODE "prompt"* — Flexible entrypoint

Examples:

bash
claude-run bedrock "explain this code"
claude-run max "optimize this function"


### Navigation Aliases

Included in .zshrc:

- cws   → ~/workspace
- ccode → ~/workspace/code
- cwhite → ~/workspace/whitepapers
- cpat  → ~/workspace/patent-pool

Modify in zshrc as needed.

---

## File Synchronization Notes

Because ~/workspace is the shared mount between *macOS* and *Lima*, the dotfiles repo automatically stays synced:

- *macOS* writes to ~/workspace/dotfiles
- *Lima* sees those same files immediately

This creates a unified environment across machines.

---

## Troubleshooting

### 1. Broken Symlinks

Recreate manually:

bash
ln -sf ~/workspace/dotfiles/zsh/zshrc ~/.zshrc
ln -sf ~/workspace/dotfiles/zsh/p10k.zsh ~/.p10k.zsh


### 2. Missing Tools

Ensure you installed:

- *brew*  
- *zsh*  
- *powerlevel10k*  
- *zsh-autosuggestions*  
- *zsh-syntax-highlighting*  

### 3. Conflicts With Existing Config

If an old .zshrc or .p10k.zsh exists, back it up:

bash
mv ~/.zshrc ~/.zshrc.backup
mv ~/.p10k.zsh ~/.p10k.zsh.backup


Then re-run the bootstrap script.

---

## License

This repository is licensed under the *MIT License*.

---

## Bootstrap Scripts

### 1. bootstrap-mac.sh

The bootstrap-mac.sh script is used to set up the macOS environment by creating symlinks to the dotfiles repository. This script is essential for the macOS setup to ensure all configuration files are linked correctly from the repo to your home directory.

### 2. bootstrap-lima.sh

The bootstrap-lima.sh script is similar to the macOS script but is tailored for the *Lima* (Linux) environment. It creates symlinks in the Lima environment, ensuring that your development environment remains consistent with the one on macOS.

Both scripts perform the following tasks:

- Create symlinks for configuration files.
- Ensure the appropriate tools (e.g., *Homebrew, **Zsh plugins*) are installed.
- Synchronize the environment between macOS and Lima.

You can run these scripts directly after cloning the repository:

For *macOS*:

bash
./bootstrap-mac.sh


For *Lima* (Linux):

bash
./bootstrap-lima.sh


These scripts make the setup fully automated and easily reproducible.

---

## Brewfile

The *Brewfile* contains a list of dependencies that can be installed via *Homebrew* on macOS and *Linuxbrew* on Lima. To install all dependencies listed in the Brewfile, you can run:

bash
brew bundle --file=~/workspace/dotfiles/Brewfile


This will ensure that all necessary tools (e.g., *Powerlevel10k, **Zsh plugins*) are installed in both environments.

---

This *README.md* is designed to guide you through setting up and maintaining your dotfiles across macOS and Lima (Linux), making it easy to configure and keep both environments synchronized.
