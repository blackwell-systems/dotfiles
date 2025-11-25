#!/usr/bin/env bash
set -euo pipefail

SESSION="$1"
SSH_DIR="$HOME/.ssh"

mkdir -p "$SSH_DIR"
chmod 700 "$SSH_DIR"

echo "🔐 Restoring SSH keys..."

ITEM=$(bw get item "SSH-Primary" --session "$SESSION")

echo "$ITEM" | jq -r '.fields[] | select(.name=="private_key").value' > "$SSH_DIR/id_ed25519"
echo "$ITEM" | jq -r '.fields[] | select(.name=="public_key").value' > "$SSH_DIR/id_ed25519.pub"

chmod 600 "$SSH_DIR/id_ed25519"
chmod 644 "$SSH_DIR/id_ed25519.pub"

echo "🔑 SSH keys restored."
