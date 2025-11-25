#!/usr/bin/env bash
set -euo pipefail

SESSION="$1"
ENV_FILE="$HOME/.local/env.secrets"

mkdir -p "$HOME/.local"

echo "🔧 Restoring environment secrets..."

if bw get item "Environment-Secrets" --session "$SESSION" >/dev/null 2>&1; then
    bw get item "Environment-Secrets" --session "$SESSION" \
        | jq -r '.fields[] | .name + "=" + .value' \
        > "$ENV_FILE"

    chmod 600 "$ENV_FILE"

    echo "export \$(grep -v '^#' $ENV_FILE | xargs)" \
        > "$HOME/.local/load-env.sh"

    chmod 700 "$HOME/.local/load-env.sh"

    echo "🔧 Environment variables restored."
else
    echo "⚠️ No Environment-Secrets vault item found."
fi

