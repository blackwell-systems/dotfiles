#!/usr/bin/env bash
set -euo pipefail

DOTFILES_DIR="$HOME/workspace/dotfiles"
VAULT_DIR="$DOTFILES_DIR/vault"
SESSION_FILE="$VAULT_DIR/.bw-session"

echo "🔐 Bitwarden bootstrap starting..."

# 1. Unlock Bitwarden ---------------------------------------------------
if [[ -f "$SESSION_FILE" ]]; then
    SESSION=$(cat "$SESSION_FILE")
    if ! bw unlock --check --session "$SESSION" >/dev/null 2>&1; then
        echo "🔓 Bitwarden session expired. Re-unlocking..."
        rm -f "$SESSION_FILE"
    fi
fi

if [[ ! -f "$SESSION_FILE" ]]; then
    SESSION=$(bw unlock --raw)
    echo "$SESSION" > "$SESSION_FILE"
    chmod 600 "$SESSION_FILE"
    echo "🔓 Vault unlocked."
else
    SESSION=$(cat "$SESSION_FILE")
    echo "🔓 Vault session restored."
fi

# 2. Restore SSH keys ---------------------------------------------------
"$VAULT_DIR/restore-ssh.sh" "$SESSION"

# 3. Restore AWS credentials & SSO profiles ------------------------------
"$VAULT_DIR/restore-aws.sh" "$SESSION"

# 4. Restore environment variables --------------------------------------
"$VAULT_DIR/restore-env.sh" "$SESSION"

echo "🎉 Bitwarden bootstrap complete."
echo "SSH keys, AWS profiles, and environment secrets are now restored."
