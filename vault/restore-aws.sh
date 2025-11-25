#!/usr/bin/env bash
set -euo pipefail

SESSION="$1"
AWS_DIR="$HOME/.aws"

mkdir -p "$AWS_DIR"
chmod 700 "$AWS_DIR"

echo "🟦 Restoring AWS configuration..."

# Fetch SSO + credentials in one Bitwarden item
ITEM=$(bw get item "AWS-Master" --session "$SESSION")

# -------------------
# Extract SSO fields
# -------------------
SSO_DEV_START_URL=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="SSO_DEV_START_URL").value')
SSO_DEV_REGION=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="SSO_DEV_REGION").value')

SSO_PROD_START_URL=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="SSO_PROD_START_URL").value')
SSO_PROD_REGION=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="SSO_PROD_REGION").value')

# -------------------
# Extract permanent access keys
# -------------------
PERSONAL_KEY=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="PERSONAL_AWS_ACCESS_KEY_ID").value')
PERSONAL_SECRET=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="PERSONAL_AWS_SECRET_ACCESS_KEY").value')

BWS_KEY=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="BWS_AWS_ACCESS_KEY_ID").value')
BWS_SECRET=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="BWS_AWS_SECRET_ACCESS_KEY").value')

# -------------------
# Extract temporary STS credentials (prod-workbench)
# -------------------
PROD_WB_KEY=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="PROD_WB_ACCESS_KEY_ID").value')
PROD_WB_SECRET=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="PROD_WB_SECRET_ACCESS_KEY").value')
PROD_WB_TOKEN=$(echo "$ITEM" | jq -r '.fields[] | select(.name=="PROD_WB_SESSION_TOKEN").value')

# -------------------
# Render AWS CONFIG FILE (sso sessions + profiles)
# -------------------
sed \
  -e "s|{{SSO_DEV_START_URL}}|$SSO_DEV_START_URL|g" \
  -e "s|{{SSO_DEV_REGION}}|$SSO_DEV_REGION|g" \
  -e "s|{{SSO_PROD_START_URL}}|$SSO_PROD_START_URL|g" \
  -e "s|{{SSO_PROD_REGION}}|$SSO_PROD_REGION|g" \
  "$HOME/workspace/dotfiles/vault/template-aws-config" \
  > "$AWS_DIR/config"

chmod 600 "$AWS_DIR/config"

echo "🟦 AWS SSO config restored."

# -------------------
# Render AWS CREDENTIALS FILE
# -------------------
sed \
  -e "s|{{PERSONAL_AWS_ACCESS_KEY_ID}}|$PERSONAL_KEY|g" \
  -e "s|{{PERSONAL_AWS_SECRET_ACCESS_KEY}}|$PERSONAL_SECRET|g" \
  -e "s|{{BWS_AWS_ACCESS_KEY_ID}}|$BWS_KEY|g" \
  -e "s|{{BWS_AWS_SECRET_ACCESS_KEY}}|$BWS_SECRET|g" \
  -e "s|{{PROD_WB_ACCESS_KEY_ID}}|$PROD_WB_KEY|g" \
  -e "s|{{PROD_WB_SECRET_ACCESS_KEY}}|$PROD_WB_SECRET|g" \
  -e "s|{{PROD_WB_SESSION_TOKEN}}|$PROD_WB_TOKEN|g" \
  "$HOME/workspace/dotfiles/vault/template-aws-credentials" \
  > "$AWS_DIR/credentials"

chmod 600 "$AWS_DIR/credentials"

echo "🔑 AWS credentials restored."
