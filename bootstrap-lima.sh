#!/usr/bin/env bash
set -euo pipefail

DOTFILES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Lima (dev-ubuntu) bootstrap starting ==="

# ------------------------------------------------------------
# 0. Basic OS packages (build tools, curl, git, zsh, etc.)
# ------------------------------------------------------------
echo "[Lima] Installing base apt packages..."
sudo apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
  build-essential \
  curl \
  file \
  git \
  zsh

# ------------------------------------------------------------
# 1. Linuxbrew (Homebrew for Linux)
# ------------------------------------------------------------
if ! command -v brew >/dev/null 2>&1; then
  echo "[Lima] Installing Linuxbrew..."
  /bin/bash -c \
    "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

  # Standard Linuxbrew shellenv
  if [ -d /home/linuxbrew/.linuxbrew/bin ]; then
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> "$HOME/.zprofile"
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> "$HOME/.profile"
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
  fi
else
  echo "[Lima] Linuxbrew already installed."
  # Make sure shellenv is loaded for this run
  if command -v brew >/dev/null 2>&1; then
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" || true
  fi
fi

# ------------------------------------------------------------
# 2. Brew Bundle (shared Brewfile)
# ------------------------------------------------------------
if command -v brew >/dev/null 2>&1; then
  if [ -f "$DOTFILES_DIR/Brewfile" ]; then
    echo "[Lima] Running brew bundle with $DOTFILES_DIR/Brewfile..."
    brew bundle --file="$DOTFILES_DIR/Brewfile"
  else
    echo "[Lima] No Brewfile found at $DOTFILES_DIR/Brewfile; skipping brew bundle."
  fi
else
  echo "[Lima] brew is not available; skipping Brewfile installation."
fi

# ------------------------------------------------------------
# 3. Sanity check: ~/workspace should be mounted from macOS
# ------------------------------------------------------------
if [ ! -d "$HOME/workspace" ]; then
  echo "ERROR: ~/workspace does not exist inside Lima."
  echo "This script assumes your macOS host is sharing ~/workspace into the VM"
  echo "via lima.yaml, e.g.:"
  echo "  mounts:"
  echo "    - location: ~/workspace"
  echo "      path: /home/ubuntu/workspace"
  echo
  echo "Fix the Lima mount, then rerun this script."
  exit 1
fi

# NOTE: We do *not* create ~/workspace or any subdirectories here.
# macOS is the source of truth for:
#   ~/workspace
#   ~/workspace/code
#   ~/workspace/whitepapers
#   ~/workspace/patent-pool
#   ~/workspace/dotfiles
# Lima just consumes that structure.

# ------------------------------------------------------------
# 4. Dotfiles symlinks (shared with macOS)
# ------------------------------------------------------------
echo "[Lima] Linking dotfiles via bootstrap-dotfiles.sh..."
"$DOTFILES_DIR/bootstrap-dotfiles.sh"

# ------------------------------------------------------------
# 5. Set default shell to zsh (if not already)
# ------------------------------------------------------------
if [ "$SHELL" != "$(command -v zsh)" ]; then
  echo "[Lima] Changing default shell to zsh..."
  chsh -s "$(command -v zsh)" || echo "WARNING: chsh failed; you may need to run it manually."
else
  echo "[Lima] Default shell already zsh."
fi

echo "=== Lima bootstrap complete ==="
echo "Next steps:"
echo "  - Open a new Lima shell so zsh + Powerlevel10k + plugins are active."
echo "  - Use aliases like: cws / ccode / cwhite / cpat once repos are cloned on macOS."
